package generator

import (
	"text/template"
)

var qsTmpl = template.Must(
	template.New("generator").
		Parse(qsCode),
)

const qsCode = `
{{ range .Configs }}
	{{ $ft := printf "%s%s" .StructName "DBField" }}
	// ===== BEGIN of query set {{ .Name }}

	// {{ .Name }} is an queryset type for {{ .StructName }}
	type {{ .Name }} struct {
		db *gorm.DB
	}

	// New{{ .Name }} constructs new {{ .Name }}
	func New{{ .Name }}(db *gorm.DB) {{ .Name }} {
		return {{ .Name }}{
			db: db.Model(&{{ .StructName }}{}),
		}
	}

	func (qs {{ .Name }}) w(db *gorm.DB) {{ .Name }} {
		return New{{ .Name }}(db)
	}

	func (qs {{ .Name }}) joinFields(fields []{{ $ft }}, prefix, suffix, separator string) string {
		names := make([]string, len(fields))
		for i, f := range fields {
			names[i] = prefix + f.String() + suffix
		}
		return strings.Join(names, separator)
	}

	func (qs {{ .Name }}) Select(fields ...{{ $ft }}) {{ .Name }} {
		return qs.w(qs.db.Select(qs.joinFields(fields, "", "", ", ")))
	}

	func (qs {{ .Name }}) OrderAscBy(fields ...{{ $ft }}) {{ .Name }} {
		return qs.w(qs.db.Order(qs.joinFields(fields, "", " ASC", ", ")))
	}

	func (qs {{ .Name }}) OrderDescBy(fields ...{{ $ft }}) {{ .Name }} {
		return qs.w(qs.db.Order(qs.joinFields(fields, "", " DESC", ", ")))
	}

	{{ range .Methods }}
		{{ .GetDoc .GetMethodName }}
		func ({{ .GetReceiverDeclaration }}) {{ .GetMethodName }}({{ .GetArgsDeclaration }})
		{{- .GetReturnValuesDeclaration }} {
		{{ .GetBody }}
		}
	{{ end }}

	// ===== END of query set {{ .Name }}

	// ===== BEGIN of {{ .StructName }} modifiers

	// {{ $ft }} describes database schema field. It requires for method 'Update'
	type {{ $ft }} string

	// String method returns string representation of field.
	func (f {{ $ft }}) String() string {
		return string(f)
	}

	// {{ .StructName }}DBSchema stores db field names of {{ .StructName }}
	var {{ .StructName }}DBSchema = struct {
		{{ range .Fields }}
			{{ .Name }} {{ $ft }}
		{{- end }}
	}{
		{{ range .Fields }}
			{{ .Name }}: {{ $ft }}("{{ .DBName }}"),
		{{- end }}
	}

	// Update updates {{ .StructName }} fields by primary key
	func (o *{{ .StructName }}) Update(db *gorm.DB, fields ...{{ $ft }}) error {
		dbNameToFieldName := map[string]interface{}{
			{{- range .Fields }}
				"{{ .DBName }}": o.{{ .Name }},
			{{- end }}
		}
		u := map[string]interface{}{}
		for _, f := range fields {
			fs := f.String()
			u[fs] = dbNameToFieldName[fs]
		}
		if err := db.Model(o).Updates(u).Error; err != nil {
			return errors.Wrapf(err, "can't update {{ .StructName }} %v fields %v", o, fields)
		}

		return nil
	}

	// {{ .StructName }}Updater is an {{ .StructName }} updates manager
	type {{ .StructName }}Updater struct {
		fields map[string]interface{}
		db *gorm.DB
	}

	// New{{ .StructName }}Updater creates new {{ .StructName }} updater
	func New{{ .StructName }}Updater(db *gorm.DB) {{ .StructName }}Updater {
		return {{ .StructName }}Updater{
			fields: map[string]interface{}{},
			db: db.Model(&{{ .StructName }}{}),
		}
	}

	// ===== END of {{ .StructName }} modifiers
{{ end }}
`
