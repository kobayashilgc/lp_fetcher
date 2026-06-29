package models

import (
	"encoding/json"
	"testing"
)

func TestInferMaterial(t *testing.T) {
	tests := []struct {
		name     string
		itemName string
		want     string
	}{
		{"CD priority", "Album CD Deluxe", "CD"},
		{"LP", "Highway To Hell LP", "黑胶"},
		{"Chinese vinyl", "爵士黑胶套装", "黑胶"},
		{"other", "Music Poster", "其他"},
		{"CD over vinyl keyword", "CD+LP Box", "CD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InferMaterial(tt.itemName); got != tt.want {
				t.Errorf("InferMaterial(%q) = %q, want %q", tt.itemName, got, tt.want)
			}
		})
	}
}

func TestMaterialFetchResultFrom(t *testing.T) {
	t.Run("groups and sorts materials", func(t *testing.T) {
		fr := FetchResult{
			Count:  3,
			Msg:    "",
			Status: "success",
			Items: []ItemBrief{
				{ItemID: "1", ItemName: "CD Album", Price: "99", Material: "CD"},
				{ItemID: "2", ItemName: "LP Album", Price: "289", Material: "黑胶"},
				{ItemID: "3", ItemName: "Poster", Price: "50", Material: "其他"},
			},
		}

		got := MaterialFetchResultFrom(fr)
		if got.Count != 3 || got.Status != "success" {
			t.Fatalf("unexpected top-level fields: %+v", got)
		}
		if len(got.Materials) != 3 {
			t.Fatalf("len(Materials) = %d, want 3", len(got.Materials))
		}
		if got.Materials[0].Material != "黑胶" || len(got.Materials[0].Items) != 1 {
			t.Fatalf("first group = %+v, want 黑胶 with 1 item", got.Materials[0])
		}
		if got.Materials[1].Material != "CD" {
			t.Fatalf("second group material = %q, want CD", got.Materials[1].Material)
		}
		if got.Materials[2].Material != "其他" {
			t.Fatalf("third group material = %q, want 其他", got.Materials[2].Material)
		}

		item := got.Materials[0].Items[0]
		if item.ItemID != "2" || item.ItemName != "LP Album" || item.Price != "289" {
			t.Fatalf("mapped item = %+v", item)
		}
	})

	t.Run("empty items yields empty array", func(t *testing.T) {
		fr := FetchResult{
			Count:  0,
			Msg:    "",
			Status: "success",
			Items:  []ItemBrief{},
		}

		got := MaterialFetchResultFrom(fr)
		if got.Materials == nil {
			t.Fatal("Materials = nil, want empty slice")
		}
		if len(got.Materials) != 0 {
			t.Fatalf("len(Materials) = %d, want 0", len(got.Materials))
		}

		data, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != `{"count":0,"msg":"","status":"success","materials":[]}` {
			t.Fatalf("json = %s", data)
		}
	})

	t.Run("failed status yields nil materials", func(t *testing.T) {
		fr := FetchResult{
			Count:  0,
			Msg:    "查询结果过多、请精确查询关键词",
			Status: "failed",
			Items:  nil,
		}

		got := MaterialFetchResultFrom(fr)
		if got.Materials != nil {
			t.Fatalf("Materials = %+v, want nil", got.Materials)
		}
	})
}

func TestNormalizeCatName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"spaces removed", "Cool Jazz 冷爵士", "cooljazz冷爵士"},
		{"already lowercase", "爵士", "爵士"},
		{"mixed case english", "ROCK Pop", "rockpop"},
		{"tabs and newlines", "New\tLine\nTest", "newlinetest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeCatName(tt.in); got != tt.want {
				t.Errorf("normalizeCatName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestLeafCateList(t *testing.T) {
	resp := CategoryResponse{
		Result: CategoryResult{
			CateList: []CategoryRaw{
				{
					CateID:   1,
					CateName: "爵士",
					ChildCateList: []CategoryRaw{
						{CateID: 2, CateName: "Cool Jazz 冷爵士"},
					},
				},
				{CateID: 3, CateName: "摇滚"},
			},
		},
	}

	got := resp.LeafCateList()
	if len(got) != 2 {
		t.Fatalf("len(LeafCateList()) = %d, want 2", len(got))
	}
	if got[0].CateID != 2 || got[0].CatName != "cooljazz冷爵士" {
		t.Fatalf("first leaf = %+v, want cateId=2 catName=cooljazz冷爵士", got[0])
	}
	if got[1].CateID != 3 || got[1].CatName != "摇滚" {
		t.Fatalf("second leaf = %+v, want cateId=3 catName=摇滚", got[1])
	}
}
