package topo_server_test

import (
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("object test", func() {
	var bizId string
	var childInstId string
	objectClient := topoServerClient.Object()
	instClient := topoServerClient.Instance()

	Describe("mainline object test", func() {
		It("create business bk_biz_name = 'abc'", func() {
			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_productor":  "",
				"bk_biz_tester":     "",
				"bk_biz_developer":  "",
				"operator":          "",
				"bk_biz_name":       "abc",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data).To(ContainElement("abc"))
			bizId = strconv.FormatInt(int64(rsp.Data["bk_biz_id"].(float64)), 10)
		})

		It("create mainline object bk_obj_id = 'test_object' and bk_obj_name='test_object'", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:     "bk_biz_topo",
					ObjectID:   "test_object",
					ObjectName: "test_object",
					OwnerID:    "0",
					ObjIcon:    "icon-cc-business",
				},
				AssociationID: "biz",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			data := map[string]interface{}{
				"bk_classification_id": "bk_biz_topo",
				"bk_obj_icon":          "icon-cc-business",
				"bk_obj_id":            "test_object",
				"bk_obj_name":          "test_object",
				"bk_supplier_account":  "0",
			}
			for k, v := range data {
				Expect(rsp.Data.(map[string]interface{})).To(HaveKeyWithValue(k, v))
			}
		})

		It("delete mainline object bk_obj_id = 'test_object'", func() {
			rsp, err := objectClient.DeleteModel(context.Background(), "0", "test_object", header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("create mainline object bk_obj_id = 'cc_test_object' and bk_obj_name='cc_test_object'", func() {
			input := &metadata.MainLineObject{
				Object: metadata.Object{
					ObjCls:     "bk_biz_topo",
					ObjectID:   "cc_test_object",
					ObjectName: "cc_test_object",
					OwnerID:    "0",
					ObjIcon:    "icon-cc-business",
				},
				AssociationID: "biz",
			}
			rsp, err := objectClient.CreateModel(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			data := map[string]interface{}{
				"bk_classification_id": "bk_biz_topo",
				"bk_obj_icon":          "icon-cc-business",
				"bk_obj_id":            "cc_test_object",
				"bk_obj_name":          "cc_test_object",
				"bk_supplier_account":  "0",
			}
			for k, v := range data {
				Expect(rsp.Data.(map[string]interface{})).To(HaveKeyWithValue(k, v))
			}
		})

		It("search mainline object", func() {
			rsp, err := objectClient.SelectModel(context.Background(), "0", header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			data1 := metadata.MainlineObjectTopo{
				ObjID:      "biz",
				ObjName:    "业务",
				OwnerID:    "0",
				NextObj:    "cc_test_object",
				NextName:   "cc_test_object",
				PreObjID:   "",
				PreObjName: "",
			}
			Expect(rsp.Data).To(ContainElement(data1))
			data2 := metadata.MainlineObjectTopo{
				ObjID:      "cc_test_object",
				ObjName:    "cc_test_object",
				OwnerID:    "0",
				NextObj:    "set",
				NextName:   "集群",
				PreObjID:   "biz",
				PreObjName: "业务",
			}
			Expect(rsp.Data).To(ContainElement(data2))
			data3 := metadata.MainlineObjectTopo{
				ObjID:      "set",
				ObjName:    "集群",
				OwnerID:    "0",
				NextObj:    "module",
				NextName:   "模块",
				PreObjID:   "cc_test_object",
				PreObjName: "cc_test_object",
			}
			Expect(rsp.Data).To(ContainElement(data3))
		})
	})

	Describe("instance topo test", func() {
		It("search instance topo", func() {
			rsp, err := objectClient.SelectInst(context.Background(), "0", bizId, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			res := []map[string]interface{}{}
			json.Unmarshal(j, &res)
			data := res[0]
			Expect(data["bk_inst_name"].(string)).To(Equal("abc"))
			Expect(data["bk_obj_id"].(string)).To(Equal("biz"))
			Expect(data["bk_obj_name"].(string)).To(Equal("业务"))
			js, err := json.Marshal(data["child"].([]interface{})[0])
			child := map[string]interface{}{}
			json.Unmarshal(js, &child)
			Expect(child["bk_inst_name"].(string)).To(Equal("cc_test_object"))
			Expect(child["bk_obj_id"].(string)).To(Equal("cc_test_object"))
			Expect(child["bk_obj_name"].(string)).To(Equal("cc_test_object"))
			Expect(len(child["child"].([]interface{}))).To(Equal(0))
			childInstId = strconv.FormatInt(int64(child["bk_inst_id"].(float64)), 10)
		})

		It("search instance topo child", func() {
			rsp, err := objectClient.SelectInstChild(context.Background(), "0", "cc_test_object", bizId, childInstId, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := []map[string]interface{}{}
			json.Unmarshal(j, &data)
			child := data[0]
			Expect(child["bk_inst_name"].(string)).To(Equal("cc_test_object"))
			Expect(child["bk_obj_id"].(string)).To(Equal("cc_test_object"))
			Expect(child["bk_obj_name"].(string)).To(Equal("cc_test_object"))
			Expect(len(child["child"].([]interface{}))).To(Equal(0))
		})

		It("search instance topo", func() {
			rsp, err := instClient.GetInternalModule(context.Background(), "0", bizId, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.SetName).To(Equal("空闲机池"))
			Expect(len(rsp.Data.Module)).To(Equal(2))
		})
	})

	Describe("classification test", func() {
		var clsId, clsId2 string

		It("create classification", func() {
			input := &metadata.Classification{
				ClassificationID:   "cc_class",
				ClassificationName: "新测试分类",
				ClassificationIcon: "icon-cc-middleware",
				OwnerID:            "0",
			}
			rsp, err := objectClient.CreateClassification(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := map[string]interface{}{}
			json.Unmarshal(j, &data)
			Expect(data["bk_classification_icon"].(string)).To(Equal("icon-cc-middleware"))
			Expect(data["bk_classification_id"].(string)).To(Equal("cc_class"))
			Expect(data["bk_classification_name"].(string)).To(Equal("新测试分类"))
			Expect(data["bk_classification_type"].(string)).To(Equal(""))
			clsId = strconv.FormatInt(int64(data["id"].(float64)), 10)
		})

		It("create classification", func() {
			input := &metadata.Classification{
				ClassificationID:   "cc_est_object",
				ClassificationName: "cc_est_object",
				ClassificationIcon: "icon-cc-default-class",
				OwnerID:            "0",
			}
			rsp, err := objectClient.CreateClassification(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := map[string]interface{}{}
			json.Unmarshal(j, &data)
			Expect(data["bk_classification_icon"].(string)).To(Equal("icon-cc-default-class"))
			Expect(data["bk_classification_id"].(string)).To(Equal("cc_est_object"))
			Expect(data["bk_classification_name"].(string)).To(Equal("cc_est_object"))
			Expect(data["bk_classification_type"].(string)).To(Equal(""))
			clsId2 = strconv.FormatInt(int64(data["id"].(float64)), 10)
		})

		It("update classification", func() {
			input := map[string]interface{}{
				"bk_classification_name": "cc模型分类",
				"bk_classification_icon": "icon-cc-default-class",
			}
			rsp, err := objectClient.UpdateClassification(context.Background(), clsId, header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("delete classification", func() {
			input := map[string]interface{}{}
			rsp, err := objectClient.DeleteClassification(context.Background(), clsId2, header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search classification", func() {
			input := map[string]interface{}{}
			rsp, err := objectClient.SelectClassificationWithParams(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			Expect(j).To(ContainSubstring("\"bk_classification_icon\":\"icon-cc-default-class\""))
			Expect(j).To(ContainSubstring("\"bk_classification_id\":\"cc_class\""))
			Expect(j).To(ContainSubstring("\"bk_classification_name\":\"cc模型分类\""))
			Expect(j).NotTo(ContainSubstring("\"bk_classification_id\":\"cc_est_object\""))
			Expect(j).NotTo(ContainSubstring("\"bk_classification_name\":\"cc_est_object\""))
			Expect(j).NotTo(ContainSubstring("\"bk_classification_name\":\"新测试分类\""))
		})
	})

	Describe("object test", func() {
		var objId string

		It("create object bk_classification_id = 'cc_class' and bk_obj_id='cc_obj'", func() {
			input := metadata.Object{
				ObjCls:     "cc_class",
				ObjIcon:    "icon-cc-business",
				ObjectID:   "cc_obj",
				ObjectName: "cc模型",
				OwnerID:    "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := metadata.Object{}
			json.Unmarshal(j, &data)
			Expect(data.ObjCls).To(Equal(input.ObjCls))
			Expect(data.ObjIcon).To(Equal(input.ObjIcon))
			Expect(data.ObjectID).To(Equal(input.ObjectID))
			Expect(data.ObjectName).To(Equal(input.ObjectName))
			Expect(data.OwnerID).To(Equal(input.OwnerID))
			Expect(data.Creator).To(Equal(input.Creator))
		})

		It("create object bk_classification_id = 'cc_class' and bk_obj_id='test_obj'", func() {
			input := metadata.Object{
				ObjCls:     "cc_class",
				ObjIcon:    "icon-cc-business",
				ObjectID:   "test_obj",
				ObjectName: "test_obj",
				OwnerID:    "0",
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := metadata.Object{}
			json.Unmarshal(j, &data)
			Expect(data.ObjCls).To(Equal(input.ObjCls))
			Expect(data.ObjIcon).To(Equal(input.ObjIcon))
			Expect(data.ObjectID).To(Equal(input.ObjectID))
			Expect(data.ObjectName).To(Equal(input.ObjectName))
			Expect(data.OwnerID).To(Equal(input.OwnerID))
			Expect(data.Creator).To(Equal(input.Creator))
			objId = strconv.FormatInt(data.ID, 10)
		})

		It("search classifications objects", func() {
			input := map[string]interface{}{
				"bk_classification_id": "cc_class",
			}
			rsp, err := objectClient.SelectClassificationWithObjects(context.Background(), "0", header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := []map[string]interface{}{}
			json.Unmarshal(j, &data)
			Expect(data[0]["bk_classification_id"].(string)).To(Equal("cc_class"))
			js, err := json.Marshal(data[0]["bk_objects"])
			objects := []metadata.Object{}
			json.Unmarshal(js, &objects)
			Expect(len(objects)).To(Equal(2))
			Expect(objects[0].ObjectID).To(Equal("cc_obj"))
			Expect(objects[1].ObjectID).To(Equal("test_obj"))
		})

		It("update object", func() {
			input := map[string]interface{}{
				"bk_obj_name": "test_obj1",
			}
			rsp, err := objectClient.UpdateObject(context.Background(), objId, header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search objects", func() {
			input := map[string]interface{}{
				"bk_obj_id":           "test_obj",
				"bk_supplier_account": "0",
			}
			rsp, err := objectClient.SelectObjectWithParams(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := []map[string]interface{}{}
			json.Unmarshal(j, &data)
			Expect(data[0]["bk_obj_id"].(string)).To(Equal("test_obj"))
			Expect(data[0]["bk_obj_name"].(string)).To(Equal("test_obj1"))
		})

		It("delete object", func() {
			input := map[string]interface{}{}
			rsp, err := objectClient.DeleteObject(context.Background(), objId, header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search object topo", func() {
			input := map[string]interface{}{
				"bk_classification_id": "cc_class",
			}
			rsp, err := objectClient.SelectObjectTopo(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search object topo graphics", func() {
			rsp, err := objectClient.SelectObjectTopoGraphics(context.Background(), "global", "0", header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			Expect(j).To(ContainSubstring("\"node_name\":\"cc模型\""))
			Expect(j).NotTo(ContainSubstring("\"node_name\":\"test_obj1\""))
		})

		It("update object topo graphics", func() {
			input := map[string]interface{}{
				"bk_obj_id":  "cc_obj",
				"bk_inst_id": 0,
				"node_type":  "obj",
				"position": map[string]interface{}{
					"x": -75,
					"y": 108,
				},
			}
			rsp, err := objectClient.UpdateObjectTopoGraphics(context.Background(), "global", "0", header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search object topo graphics", func() {
			rsp, err := objectClient.SelectObjectTopoGraphics(context.Background(), "global", "0", header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			Expect(j).NotTo(ContainSubstring("\"position\":{\"x\":-75,\"y\":108}"))
		})
	})

	Describe("object attribute group test", func() {
		var groupId int64

		Describe("group test", func() {
			var group metadata.Group

			It("create group bk_obj_id='cc_obj'", func() {
				input := metadata.Group{
					GroupID:    "1",
					GroupName:  "123",
					GroupIndex: 1,
					ObjectID:   "cc_obj",
					OwnerID:    "0",
				}
				rsp, err := objectClient.CreatePropertyGroup(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.Group{}
				json.Unmarshal(j, &data)
				Expect(data.GroupID).To(Equal(input.GroupID))
				Expect(data.GroupName).To(Equal(input.GroupName))
				Expect(data.GroupIndex).To(Equal(input.GroupIndex))
				Expect(data.ObjectID).To(Equal(input.ObjectID))
				Expect(data.OwnerID).To(Equal(input.OwnerID))
				groupId = data.ID
			})

			It("delete group", func() {
				rsp, err := objectClient.DeletePropertyGroup(context.Background(), strconv.FormatInt(groupId, 10), header)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("create group bk_obj_id='cc_obj'", func() {
				input := metadata.Group{
					GroupID:    "1",
					GroupName:  "123",
					GroupIndex: 1,
					ObjectID:   "cc_obj",
					OwnerID:    "0",
				}
				rsp, err := objectClient.CreatePropertyGroup(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.Group{}
				json.Unmarshal(j, &data)
				Expect(data.GroupID).To(Equal(input.GroupID))
				Expect(data.GroupName).To(Equal(input.GroupName))
				Expect(data.GroupIndex).To(Equal(input.GroupIndex))
				Expect(data.ObjectID).To(Equal(input.ObjectID))
				Expect(data.OwnerID).To(Equal(input.OwnerID))
				groupId = data.ID
				group = data
			})

			It("update group", func() {
				input := &metadata.PropertyGroupCondition{
					Condition: map[string]interface{}{
						"id": groupId,
					},
					Data: map[string]interface{}{
						"bk_group_name": "456",
					},
				}
				rsp, err := objectClient.UpdatePropertyGroup(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				group.GroupName = "456"
			})

			It("search group bk_obj_id='cc_obj'", func() {
				input := map[string]interface{}{}
				rsp, err := objectClient.SelectPropertyGroupByObjectID(context.Background(), "0", "cc_obj", header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := []metadata.Group{}
				json.Unmarshal(j, &data)
				Expect(data).To(ContainElement(group))
			})
		})

		Describe("object attribute test", func() {
			var attrId, attrId1 string

			It("create object attribute bk_obj_id='cc_obj' and bk_property_id='test_sglchar' and bk_property_name='test_sglchar'", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "cc_obj",
						PropertyID:    "test_sglchar",
						PropertyName:  "test_sglchar",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "singlechar",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.Attribute{}
				json.Unmarshal(j, &data)
				Expect(data.ObjectID).To(Equal(input.ObjectID))
				Expect(data.PropertyID).To(Equal(input.PropertyID))
				Expect(data.PropertyName).To(Equal(input.PropertyName))
				Expect(data.PropertyGroup).To(Equal(input.PropertyGroup))
				Expect(data.IsEditable).To(Equal(input.IsEditable))
				Expect(data.PropertyType).To(Equal(input.PropertyType))
				Expect(data.OwnerID).To(Equal(input.OwnerID))
				attrId = strconv.FormatInt(data.ID, 10)
			})

			It("create object attribute bk_obj_id='cc_obj' and bk_property_id='test_sglchar' and bk_property_name='test_sglchar'", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "cc_obj",
						PropertyID:    "test_singlechar",
						PropertyName:  "test_singlechar",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "singlechar",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := metadata.Attribute{}
				json.Unmarshal(j, &data)
				Expect(data.ObjectID).To(Equal(input.ObjectID))
				Expect(data.PropertyID).To(Equal(input.PropertyID))
				Expect(data.PropertyName).To(Equal(input.PropertyName))
				Expect(data.PropertyGroup).To(Equal(input.PropertyGroup))
				Expect(data.IsEditable).To(Equal(input.IsEditable))
				Expect(data.PropertyType).To(Equal(input.PropertyType))
				Expect(data.OwnerID).To(Equal(input.OwnerID))
				attrId1 = strconv.FormatInt(data.ID, 10)
			})

			It("update object attribute id="+attrId1, func() {
				input := map[string]interface{}{
					"bk_property_name": "ayayyaya",
				}
				rsp, err := apiServerClient.UpdateObjectAtt(context.Background(), attrId1, header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("delete object attribute id="+attrId, func() {
				rsp, err := apiServerClient.DeleteObjectAtt(context.Background(), attrId, header)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("search object attribute", func() {
				input := mapstr.MapStr{
					"bk_obj_id": "cc_obj",
				}
				rsp, err := apiServerClient.GetObjectAttr(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := []map[string]interface{}{}
				json.Unmarshal(j, &data)
				Expect(data).To(ContainElement(HaveKeyWithValue("bk_property_name", "ayayyaya")))
				Expect(data).NotTo(ContainElement(HaveKeyWithValue("bk_property_name", "test_singlechar")))
				Expect(data).NotTo(ContainElement(HaveKeyWithValue("bk_property_name", "test_sglchar")))
			})
		})

		Describe("object attribute group test", func() {
			It("update property group", func() {
				arr := []metadata.PropertyGroupObjectAtt{
					metadata.PropertyGroupObjectAtt{},
				}
				arr[0].Condition.ObjectID = "cc_obj"
				arr[0].Condition.PropertyID = "test_singlechar"
				arr[0].Condition.OwnerID = "0"
				arr[0].Data.PropertyGroupID = "1" // strconv.FormatInt(groupId, 10)
				input := map[string]interface{}{
					"origin": arr,
				}
				rsp, err := objectClient.UpdatePropertyGroupObjectAtt(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("search object attribute", func() {
				input := mapstr.MapStr{
					"bk_obj_id": "cc_obj",
				}
				rsp, err := apiServerClient.GetObjectAttr(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := []map[string]interface{}{}
				json.Unmarshal(j, &data)
				Expect(data).To(ContainElement(HaveKeyWithValue("bk_property_group_name", "123")))
				Expect(data).To(ContainElement(HaveKeyWithValue("bk_property_group", "1")))
			})

			It("delete object attribute property group", func() {
				rsp, err := objectClient.DeletePropertyGroupObjectAtt(context.Background(), "0", "cc_obj", "test_singlechar", strconv.FormatInt(groupId, 10), header)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("search object attribute", func() {
				input := mapstr.MapStr{
					"bk_obj_id": "cc_obj",
				}
				rsp, err := apiServerClient.GetObjectAttr(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				j, err := json.Marshal(rsp.Data)
				data := []map[string]interface{}{}
				json.Unmarshal(j, &data)
				Expect(data).NotTo(ContainElement(HaveKeyWithValue("bk_property_group_name", "ayayyaya")))
				Expect(data).NotTo(ContainElement(HaveKeyWithValue("bk_property_group", "test_singlechar")))
			})
		})
	})

	Describe("set test", func() {
		var setId, setId1 string
		childInstIdInt, _ := strconv.ParseInt(childInstId, 10, 64)
		bizIdInt, _ := strconv.ParseInt(bizId, 10, 64)

		It("create set bk_biz_id="+bizId+" and bk_parent_id="+childInstId, func() {
			input := mapstr.MapStr{
				"bk_set_name":         "cc_set",
				"bk_parent_id":        childInstIdInt,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizIdInt,
				"bk_service_status":   "1",
				"bk_set_env":          "3",
			}
			rsp, err := instClient.CreateSet(context.Background(), bizId, header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			for k, v := range input {
				Expect(rsp.Data.(map[string]interface{})).To(HaveKeyWithValue(k, v))
			}
			setId = strconv.FormatInt(rsp.Data["bk_set_id"].(int64), 10)
		})

		It("create set bk_biz_id="+bizId+" and bk_parent_id="+childInstId, func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test",
				"bk_parent_id":        childInstIdInt,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizIdInt,
				"bk_service_status":   "1",
				"bk_set_env":          "3",
			}
			rsp, err := instClient.CreateSet(context.Background(), bizId, header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			for k, v := range input {
				Expect(rsp.Data.(map[string]interface{})).To(HaveKeyWithValue(k, v))
			}
			setId1 = strconv.FormatInt(rsp.Data["bk_set_id"].(int64), 10)
			fmt.Println(setId1)
		})

		It("update set", func() {
			input := map[string]interface{}{}
			rsp, err := instClient.UpdateSet(context.Background(), bizId, setId, header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})
	})
})
