package crud_generator

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getType(strct interface{}) string {
	t := reflect.TypeOf(strct)
	return t.Name()
}

func getFields(strct interface{}) (fields []string) {
	// Returns fields names (without ID) in lowerCase
	val := reflect.Indirect(reflect.ValueOf(strct))
	fieldNum := val.Type().NumField()
	for i := 0; i < fieldNum; i++ {
		fieldName := toSnakeCase(val.Type().Field(i).Name)
		if fieldName != "id" {
			fields = append(fields, fieldName)
		}
	}
	return
}

var CRUDCode = `// %[1]s functions
	func (%[7]sDaoImpl) Create%[1]s(%[2]s *models.%[1]s) (err error) {
	res, err := Db.NamedExec("INSERT INTO %[3]s (%[4]s) VALUES (%[5]s)", &%[2]s)
	if err != nil {
	return
	} else {
	id, err := res.LastInsertId()
	if err != nil {
	return err
	} else {
	%[2]s.ID = uint(id)
	}
	}
	return
	}

	func (%[7]sDaoImpl) Update%[1]s(%[2]s *models.%[1]s) (err error){
	_, err = Db.NamedExec("update %[3]s set %[8]s where id =:id", &%[2]s)
	return
	}

	func (%[7]sDaoImpl) Delete%[1]sByID(id uint) (err error) {
	_, err = Db.Exec("DELETE FROM %[3]s WHERE id=?", id)
	return
	}

	func (%[7]sDaoImpl) Get%[1]sByID(id uint) (%[2]s *models.%[1]s, err error) {
	err = Db.QueryRowx("SELECT * FROM %[3]s WHERE id=?", id).StructScan(&%[2]s)
	return
	}

	func (%[7]sDaoImpl) %[1]sList() (%[2]sList []models.%[1]s, err error) {
	rows, err := Db.Queryx("SELECT * FROM %[3]s")
	if err != nil {
	return
	}
	for rows.Next() {
	var %[2]s models.%[1]s
	err = rows.StructScan(&%[2]s)
	if err != nil {
	return
	}
	%[2]sList = append(%[2]sList, %[2]s)
	}
	rows.Close()
	return
	}

	func (%[7]sDaoImpl) DeleteAll%[1]ss() (err error) {
		_, err = Db.Exec("DELETE FROM %[3]s")
		return
	}
	`
var CRUDInterfaceCode = `// %[1]s functions
	Create%[1]s(%[2]s *models.%[1]s) (err error)
	Update%[1]s(%[2]s *models.%[1]s) (err error)
	Delete%[1]sByID(id uint) (err error)
	Get%[1]sByID(id uint) (%[2]s *models.%[1]s, err error)
	%[1]sList() (%[2]sList []models.%[1]s, err error)
	DeleteAll%[1]ss() (err error)
	`

func generateCRUD(name string, modelStruct interface{}) (code string) {
	modelName := getType(modelStruct)
	fields := getFields(modelStruct)

	lcModelName := string(unicode.ToLower(rune(modelName[0]))) + modelName[1:]
	scModelName := toSnakeCase(modelName)
	params1 := strings.Join(fields, ", ")
	params2 := ":" + strings.Join(fields, ", :")
	params3 := strings.Join(fields, "=?, ") + "=?"
	var params4 string
	for i, field := range fields {
		if i == len(field)-1 {
			params4 += field + "=:" + field
		} else {
			params4 += field + "=:" + field + ", "
		}
	}
	code = fmt.Sprintf(CRUDCode, modelName, lcModelName, scModelName, params1, params2, params3, name, params4)

	return
}

func generateCRUDInterface(modelStruct interface{}) (code string) {
	modelName := getType(modelStruct)

	lcModelName := string(unicode.ToLower(rune(modelName[0]))) + modelName[1:]
	scModelName := toSnakeCase(modelName)

	code = fmt.Sprintf(CRUDInterfaceCode, modelName, lcModelName, scModelName)
	return
}

func generateCRUDs(name string, modelStructs []interface{}) (code string) {
	for _, modelStruct := range modelStructs {
		code += generateCRUD(name, modelStruct)
	}
	return
}

func generateCRUDInterfaces(modelStructs []interface{}) (code string) {
	for _, modelStruct := range modelStructs {
		code += generateCRUDInterface(modelStruct)
	}
	return
}

func GenerateFiles(name string, modelStructs ...interface{}) {
	cruds := generateCRUDs(name, modelStructs)

	err := ioutil.WriteFile(name, []byte(cruds), 0644)
	if err != nil {
		panic(err)
	}

	interfaces := generateCRUDInterfaces(modelStructs)

	err = ioutil.WriteFile(name+"dao", []byte(interfaces), 0644)
	if err != nil {
		panic(err)
	}
}
