package spider

import (
	"testing"
	"time"
)

var xmlStr = `
<Spider>
    <Name>解析xml测试</Name>

    <Init>
    <Script param="aid, dataList">
        url = "http://zhihu.sogou.com/zhihu?query=golang+logo";
        method = "GET";
        req = aid.NewRequest(url, method);
        dataList.Push(req);
    </Script>
    </Init>

    <MaxDepth>1</MaxDepth>

    <AcceptedPrimaryDomain>
    <PrimaryDomain>zhihu.com</PrimaryDomain>
    </AcceptedPrimaryDomain>
    
    <DataArgs>
        <Request>
            <BufferCap>50</BufferCap>
            <MaxBufferNumber>1000</MaxBufferNumber>
        </Request>

        <Response>
            <BufferCap>50</BufferCap>
            <MaxBufferNumber>10</MaxBufferNumber>
        </Response>

        <Item>
            <BufferCap>50</BufferCap>
            <MaxBufferNumber>100</MaxBufferNumber>
        </Item>

        <Error>
            <BufferCap>50</BufferCap>
            <MaxBufferNumber>1</MaxBufferNumber>
        </Error>
    </DataArgs>

    <ModuleArgs>
        <DownloaderNumber>1</DownloaderNumber>
        <AnalyzerNumber>1</AnalyzerNumber>
        <PipelineNumber>1</PipelineNumber>
    </ModuleArgs>

    <ResponseParser>
        <Script param="aid,dataList,errorList,resp">
            var query = resp.GetDom();
            dom = query[0];
            var adom = dom.Find("a");
            var alen = adom.Length();
            var i = 0;
            for(i = 0; i != alen; i ++){
                var nowhref = adom.Slice(i,i+1).Attr("href");
                href = nowhref[0];
                if (!nowhref[1] || href == "" || href == "#" || href == "/") {
                    continue;
                }
                trimhref = nowhref[0].trim();
                lowhref = trimhref.toLowerCase();
                var fdStart = lowhref.indexOf("javascript")
                if(fdStart == 0) {
                    continue;
                }
                reqURL = resp.HTTPResp().Request.URL;
                aURLAndError = aid.ParseURL(lowhref);
                aURL = aURLAndError[0];
                if (!aURL.IsAbs()) {
                    aURL = reqURL.ResolveReference(aURL);
                }
                req = aid.NewRequest(aURL.String(), "GET");
                dataList.Push(req);
                item = aid.NewItem({Url: aURL.String()});
                dataList.Push(item);
            }
        </Script>
    </ResponseParser>

    <ItemProcessor>
        <Script param="aid, item, errorList">
        </Script>
    </ItemProcessor>
</Spider>
`

func TestParseJS(t *testing.T) {
	spider, err := GetSpiderByXML(xmlStr)
	if err != nil {
		t.Fatalf("get spider by xml fail: %s", err)
	}
	sched, err := spider.GenAndStartScheduler()
	if err != nil {
		t.Fatalf("gen and start schduler failt: %s", err)
	}
	go func() {
		for {
			err := <-sched.ErrorChan()
			if err != nil {
				logger.Error(err)
			} else {
				break
			}
		}
	}()
	for i := 0; i < 5; i++ {
		if !sched.Idle() {
			i = 0
		}
		time.Sleep(1 * time.Second)
	}
	sched.Stop()
	time.Sleep(1 * time.Second)
}
