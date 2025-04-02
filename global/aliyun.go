package global

import (
	bailian20231229 "github.com/alibabacloud-go/bailian-20231229/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

var (
	ALiYunClient             *bailian20231229.Client
	ALiYunRuntime            = &util.RuntimeOptions{}
	AliYunHeaders            = make(map[string]*string)
	ALiYunBaiLianCateID      = tea.String("cate_282029dd140c4120a310dc7ce5291636_11524451")
	ALiYunBaiLianWorkspaceID = tea.String("llm-b2aycxo90h82spog")
	ALiYunIndexID            = tea.String("novtr4vp8n")
)
