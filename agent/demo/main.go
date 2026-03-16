package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino-ext/adk/backend/local"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/middlewares/filesystem"
	"github.com/cloudwego/eino/adk/middlewares/skill"
	"github.com/cloudwego/eino/schema"
)

//func main() {
//	ctx := context.Background()
//
//	apiKey := "sk-dc5b112e28af4af5b125717ec07d72f9"
//	modelName := "qwen-plus"
//	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
//		BaseURL:     "https://dashscope.aliyuncs.com/compatible-mode/v1",
//		APIKey:      apiKey,
//		Timeout:     0,
//		Model:       modelName,
//		MaxTokens:   of(2048),
//		Temperature: of(float32(0.7)),
//		TopP:        of(float32(0.7)),
//	})
//
//	if err != nil {
//		fmt.Printf("NewChatModel of qwen failed, err=%v", err)
//	}
//
//	resp, err := chatModel.Stream(ctx, []*schema.Message{
//		schema.UserMessage("你好?"),
//	})
//	if err != nil {
//		fmt.Printf("Generate of qwen failed, err=%v", err)
//	}
//
//	//fmt.Printf("output: \n%v", resp)
//	defer resp.Close()
//
//	i := 0
//	for {
//		message, err := resp.Recv()
//		if err == io.EOF {
//			return
//		}
//		if err != nil {
//			fmt.Printf("recv failed: %v", err)
//		}
//		fmt.Print(message)
//		i++
//	}
//
//}
//
//func of[T any](v T) *T {
//	return &v
//}

//func main() {
//	ctx := context.Background()
//
//	// 初始化模型
//	model := model()
//
//	// 创建 ChatModelAgent
//	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
//		Name:        "hello_agent",
//		Description: "A friendly greeting assistant",
//		Instruction: "You are a friendly assistant. Please respond to the user in a warm tone.",
//		Model:       model,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 创建 Runner
//	runner := adk.NewRunner(ctx, adk.RunnerConfig{
//		Agent:           agent,
//		EnableStreaming: false,
//	})
//
//	// 执行对话
//	input := []adk.Message{
//		schema.UserMessage("你好，请介绍一下你自己."),
//	}
//
//	events := runner.Run(ctx, input)
//	for {
//		event, ok := events.Next()
//		if !ok {
//			break
//		}
//
//		if event.Err != nil {
//			log.Printf("错误: %v", event.Err)
//			break
//		}
//
//		fmt.Printf("Agent: %+v\n", event.Output)
//	}
//}

func model() *qwen.ChatModel {
	baseUrl := "https://dashscope.aliyuncs.com/compatible-mode/v1"
	apiKey := "sk-dc5b112e28af4af5b125717ec07d72f9"
	modelName := "qwen3.5-flash"

	// 初始化模型
	model, err := qwen.NewChatModel(context.Background(), &qwen.ChatModelConfig{
		APIKey:  apiKey,
		Model:   modelName,
		BaseURL: baseUrl,
	})
	if err != nil {
		log.Fatal(err)
	}
	return model
}

func main() {
	ctx := context.Background()

	chatModel := model()

	backend, _ := local.NewBackend(ctx, &local.Config{})

	fromFilesystem, _ := skill.NewBackendFromFilesystem(ctx, &skill.BackendFromFilesystemConfig{
		Backend: backend,
		BaseDir: "skills",
	})
	skillMiddleware, _ := skill.NewMiddleware(ctx, &skill.Config{
		Backend: fromFilesystem,
	})
	fsm, _ := filesystem.New(ctx, &filesystem.MiddlewareConfig{
		Backend: backend,
	})
	agent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "角色提取专家",
		Description: "你是一个角色提取专家，你需要使用合适的技能来提取角色",
		Model:       chatModel,
		Handlers:    []adk.ChatModelAgentMiddleware{fsm, skillMiddleware},
	})
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	content := "没有人料到，2026年的除夕夜将是人类历史上最后一个春节。随着一场史无前例的极端寒潮降临，自此世界再无春日。回家过年的我因大雪被困在老家大山里，外面温度骤降至零下 80 度，电力系统崩溃，物资匮乏，村里的亲戚打着“亲人”的名义，先是抢光了我们家的煤炭和粮食，最后为了取暖，竟然把我们一家三口扔到了零下八十度的雪地里！我们一家人被赶出去，身上的羽绒服都被抢走，刺骨的寒意像无数把钢刀刮着骨头，我眼睁睁看着二婶那张狰狞扭曲的脸在眼前放大。她手里死死拽着我妈身上仅剩的一件破棉袄，一脚踹在我已经僵硬的肚子上。旁边那个三百斤的堂弟林大宝，正骑在我爸冻僵的尸体上，疯狂地啃食着我爸手里那半个冻硬的馒头。我想喊，喉咙里却只发出破风箱般的嘶鸣。太痛了！我们一家人竟然活生生冻死在除夕夜。重来一回，我发现我正躺在温暖的床上。那种血液凝固、肢体坏死的剧痛，让我的灵魂都在颤抖。我猛地坐起身，大口大口地喘着粗气，冷汗瞬间浸透了睡衣。眼前不是漫天的风雪，也不是尸横遍野的雪地。是熟悉的粉色窗帘，是老家那盏昏黄的吊灯，房间内是喜庆的过年装扮，还有窗外传来的热闹的鞭炮响。\n我颤抖着手抓过床头的手机。2026年2月13日，腊月二十六。距离那场毁灭全人类的极寒末世，还有四天。我活过来了？“晚晚，怎么了？做噩梦了？”卧室门被猛地推开，老妈系着围裙，一脸焦急地冲进来。紧接着是老爸，手里还拿着锅铲。看着二老那红润有肉的脸庞，而不是前世那两具青紫干瘪的尸体，我的眼泪瞬间决堤。\n“爸！妈！”我扑进老妈怀里，嚎啕大哭。那种失而复得的狂喜，还有前世惨死的滔天恨意，冲击着我的天灵盖。“这孩子，怎么哭成这样？”老爸手足无措地给我擦泪，“是不是二婶又在外面骂你了？”听到“二婶”两个字，我浑身的肌肉瞬间紧绷，上一世，就是这家人，在极寒降临后，先是抢光了我们家的煤炭和粮食，最后为了取暖，竟然把我们一家三口扔到了零下八十度的雪地里！“林晚！你个懒猪，几点了还不起来干活？”门外，传来了一道尖锐刺耳的公鸭嗓。是二婶刘翠芬。哪怕隔着一道门，我也能闻到她身上那股令人作呕的贪婪味。砸门声震天响。“大哥！大嫂！你们家这闺女也太娇气了吧？大宝都饿了，怎么还没做饭？那半扇猪肉呢？赶紧拿出来给大宝炖了！”老妈的身子僵了一下，下意识地就要往外走：“哎，来了，他二婶你别急⋯⋯”我一把死死拽住老妈的手腕。指甲几乎嵌进肉里。“别去。”我抬起头，死死盯着老妈的眼睛，咬牙切齿：“妈，如果你还记得我们是怎么死的，就别去。”老妈愣住了。老爸手里的锅铲“哐当”一声掉在地上。三人对视，空气仿佛凝固。我在他们眼中，看到了同样的恐惧、震惊，以及⋯⋯觉醒后的决绝。原来，回来的不止我一个。"
	systemPrompt := `分析内容，并提取出相应的角色，
		提取规则：
		- 提取角色姓名
		- 提取角色身份
		- 提取角色描述，需要包含：体态、面部、穿着等视觉特征（0动作细节）
		使用以下结构返回

		[
		  {
			"role_name": "角色姓名",
			"identity": "角色身份",
			"personality": "角色性格",
			"role_description": "角色描述"
		  }
		]
		
		要求：
		1. 除了json结构，不允许返回其他任何内容，务必保证json结构100%正确
		2. json结构绝不使用任何markdown语法进行包裹
	`

	messages := []adk.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(content),
	}
	iterator := runner.Run(ctx, messages)
	for {
		event, ok := iterator.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Printf("Error: %v\n", event.Err)
			break
		}
		if msg, err := event.Output.MessageOutput.GetMessage(); err == nil {
			fmt.Printf("Agent: %s\n", msg.Content)
		}

		//prints.Event(event)
	}

	//list, err := fromFilesystem.List(ctx)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//for _, v := range list {
	//	fmt.Printf("skill item from file: %v\n", v)
	//}

}
