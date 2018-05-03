package main

import (
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
	req := command.Read()
	files := req.GetProtoFile()

	vanity.ForEachFieldInFiles(files, SetBsonTagFieldOption)

	resp := command.Generate(req)
	command.Write(resp)
}

// SetBsonTagFieldOption - applies bson tags
func SetBsonTagFieldOption(field *descriptor.FieldDescriptorProto) {
	moreTagsExtension := gogoproto.GetMoreTags(field)

	var moreTags reflect.StructTag
	if moreTagsExtension != nil {
		moreTags = reflect.StructTag(*moreTagsExtension)
	}

	// The value we'll be assigning to moreTags
	// start with any existing tags
	value := string(moreTags)
	if _, ok := moreTags.Lookup("bson"); !ok {
		bsonName := *field.Name
		bsonTag := bsonName + ",omitempty"
		repeatedNativeType := (!field.IsMessage() && !gogoproto.IsCustomType(field) && field.IsRepeated())
		if !gogoproto.IsNullable(field) && !repeatedNativeType {
			bsonTag = bsonName
		}

		if value != "" {
			value += " "
		}
		value += fmt.Sprintf(`bson:"%s"`, bsonTag)
	}

	if field.Options == nil {
		field.Options = &descriptor.FieldOptions{}
	}
	if err := proto.SetExtension(field.Options, gogoproto.E_Moretags, &value); err != nil {
		panic(err)
	}
}
