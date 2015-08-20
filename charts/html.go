package charts

import (
	"bytes"
	"log"
	"html/template"
)

var funcs = template.FuncMap{
	"add": func(a int, b int) int { return a + b; },
	"mul": func(a int, b int) int { return a * b; },
}

var CHART_TEMPLATE = template.Must(template.New("chart").Funcs(funcs).Parse(`
<div class="chart-info">
	{{.Meta}}
</div>
<div class="chart">
	<div class="chart-fields">
		<div class="chart-field">region</div>
		{{range $col := (index .Rows 0).Images}}
			<div class="chart-field">
				{{(index $col.Meta "stain")}}
			</div>
		{{end}}
		<div class="chart-field">overlap</div>
	</div>
	{{range $row := .Rows}}
		<div class="chart-row">
			<div class="chart-row-name">
				{{$row.Meta}}
			</div>
			{{range $col := $row.Images}}
				<div class="chart-img">
					<img src="file:///{{$col.Path}}"/>
				</div>
			{{end}}
			<div class="chart-img chart-img-overflow">
			{{range $idx, $col := $row.Images}}
				<img class="chart-overlap{{if ne $idx 0}} chart-overlap-others{{end}}"
					src="file:///{{$col.Path}}" style="left: -{{add (mul 250 $idx) (mul 4 $idx)}}px;"/>
			{{end}}
			</div>
		</div>
	{{end}}
</div>
`))

var CHARTS_TEMPLATE = template.Must(template.New("charts").Parse(`<!DOCTYPE html>
<html>
<head>
<style>
.chart {
	display: table
}
.chart-fields {
	display: table-header-group;
}
.chart-field {
	display: table-cell;
}
.chart-row {
	display: table-row;
}
.chart-row-name {
	display: table-cell;
	vertical-align: middle;
}
.chart-img {
	width: 250px;
	height: 250px;
	display: table-cell;
}
.chart-img-overflow {
	width: 800px;
	display: table-cell;
}
.chart-img img {
	width: inherit;
	height: inherit;
}
img.chart-overlap {
	position: relative;
	width: 250px;
	top: inherit;
}
img.chart-overlap-others {
	opacity: .25;
}
</style>
</head>
<body>
{{range $chart := .charts}}
{{$chart.HTML}}
<hr/>
{{end}}
</body>
</html>
`))

func ChartsHTML(charts []*Chart) template.HTML {
	return html(CHARTS_TEMPLATE, map[string]interface{}{
		"charts": charts,
	})
}

func (c *Chart) HTML() template.HTML {
	return html(CHART_TEMPLATE, c)
}

func html(t *template.Template, data interface{}) template.HTML {
	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		log.Panic(err)
	}
	return template.HTML(buf.String())
}

