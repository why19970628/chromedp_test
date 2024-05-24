package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/chromedp/chromedp"
	pkgErr "github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	html := getHtml("您好", "ninhao")
	base64, err := Html2Image(context.Background(), html)
	fmt.Println("base64", base64)
	fmt.Println("err", err)
}

func getHtml(ask, answer string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>我与汽车大师GPT的对话</title>
    <style>
        * {
            margin: 0;
            padding: 0;
        }

        html {
            background-color: #f7f7f7;
        }

        .container {
            position: relative;
            width: 318px;
            top: 130px;
            left: 0;
            right: 0;
            margin: auto;
            background: linear-gradient(180deg, #ffffff 0%%, #f6f7fb 15%%, #f6f7fb 86%%, #ffffff 100%%);
            border-radius: 10px;
            overflow: hidden;
        }

        .header {
            display: block;
            height: 87px;
            padding-top: 26px;
            background: url("https://oss.qcds.com/assets/image/web/gpt_bg.png") no-repeat top;
            background-size: cover;
            background-color: #fff;
        }

        .title {
            margin-top: 0;
            display: flex;
            height: fit-content;
            align-items: center;
            justify-content: center;
            font-family: PingFangSC, PingFang SC;
            font-weight: bold;
            font-size: 16px;
            color: #333333;
            line-height: 22px;
            font-style: normal;
        }

        .title img {
            margin-right: 4px;
            width: 21px;
            height: 21px;
        }

        .tip {
            margin-top: 9px;
            font-family: PingFangSC, PingFang SC;
            font-weight: 400;
            font-size: 10px;
            color: #999999;
            line-height: 14px;
            text-align: center;
            font-style: normal;
        }

        .talk {
            padding: 0 10px;
            box-sizing: border-box;
            font-family: PingFangSC, PingFang SC;
            font-weight: 400;
            font-size: 12px;
            color: #333333;
            line-height: 17px;
            font-style: normal;
        }

        .talk-answer {
            margin-top: 15px;
            max-width: 298px;
            background: #ffffff;
            padding: 10px;
            box-sizing: border-box;
            border-radius: 10px 10px 10px 2px;
        }

        .talk-ask {
            margin-left: auto;
            margin-top: 15px;
            max-width: 248px;
            width: fit-content;
            background: #fcded5;
            padding: 10px;
            box-sizing: border-box;
            border-radius: 10px 10px 2px 10px;
        }

        .footer {
            position: relative;
            padding-top: 24px;
            padding-bottom: 22px;
        }

        .footer p {
            margin-left: 25px;
            font-family: PingFangSC, PingFang SC;
            font-weight: 500;
            font-size: 12px;
            color: #333333;
            line-height: 17px;
            text-align: left;
            font-style: normal;
        }

        .footer img {
            position: absolute;
            right: 25px;
            top: 20px;
        }
    </style>
</head>
<body>
<div class="container">
    <div class="header">
        <div class="title">
            <img src="https://oss.qcds.com/assets/image/web/icon_logo.png" alt="" />我与汽车大师GPT的对话
        </div>
        <div class="tip">该内容由AI生成，可能出现错误或意外情况</div>
    </div>
    <div class="talk">
        	<div class="talk-ask">%s</div>
			<div id="output" class="talk-answer">%s</div>
    </div>
    <div class="footer">
        <p>立即扫码</p>
        <p>体验更多汽车智能服务～</p>
        <img src="https://oss.qcds.com/assets/image/web/icon_logo.png" alt="" />
    </div>
</div>
</body>
</html>
`, ask, answer)
}

func Html2Image(ctx context.Context, htmlContent string) (base64Str string, err error) {
	// 创建临时 HTML 文件
	tmpfile, err := ioutil.TempFile("/tmp", "temp-*.html")

	if err != nil {
		log.Println("RpcShareImage.dpCtx err 0:", err)
		return base64Str, pkgErr.WithMessage(err, "ioutil.TempFile() err")
	}
	defer os.Remove(tmpfile.Name()) // 确保在程序结束时删除临时文件

	defer func() {
		if err := tmpfile.Close(); err != nil {
			fmt.Println("tmpfile.Close() err", err)
		}
	}()

	log.Println("RpcShareImage.dpCtx err 1:", err)

	// 写入 HTML 内容到临时文件
	if _, err := tmpfile.Write([]byte(htmlContent)); err != nil {
		return base64Str, pkgErr.WithMessage(err, "tmpfile.Write() err")
	}

	log.Println("RpcShareImage.dpCtx err 2:", err)

	// 生成 chrome 实例
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	//defer func() {
	//	if err := chromedp.Cancel(dpCtx); err != nil {
	//		log.Println("RpcShareImage.dpCtx err 4:", err)
	//	}
	//}()
	defer cancel()

	// 设置截图某一个元素
	var buf []byte
	filePath := tmpfile.Name()

	if err := chromedp.Run(ctx, elementScreenshot(fmt.Sprintf(`file:%s`, filePath), `div.container`, &buf)); err != nil {
		return base64Str, pkgErr.WithMessage(err, "chromedp.Run() err")
	}

	log.Println("RpcShareImage.Html2Image 3:", err)

	if err := ioutil.WriteFile("./h3.png", buf, 0o644); err != nil {
		return base64Str, pkgErr.WithMessage(err, "ioutil.WriteFile() err")
	}
	base64Str = base64.StdEncoding.EncodeToString(buf)
	return base64Str, nil
}

// elementScreenshot 截图页面某一个元素
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
