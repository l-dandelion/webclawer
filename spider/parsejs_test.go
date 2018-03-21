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
        url = "http://codeforces.com/api/user.info?handles=DmitriyH;Fefer_Ivan";
        method = "GET";
        req = aid.NewRequest(url, method);
        dataList.Push(req);
    </Script>
    </Init>

    <MaxDepth>0</MaxDepth>

    <AcceptedPrimaryDomain>
    <PrimaryDomain>codeforce.com</PrimaryDomain>
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
        <Script param="aid, dataList, errorList, resp">
            val = resp.GetText();

            function bin2String(array) {
              return String.fromCharCode.apply(String, array);
            }

            jsonObj = JSON.parse(bin2String(val[0]));
            console.log(jsonObj["status"])

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
				t.Log(err)
			} else {
				break
			}
		}
	}()
	for i := 0; i < 5; i++ {
		if !sched.Idle() {
			i = 0
		}
		time.Sleep(100 * time.Millisecond)
	}
	sched.Stop()
	time.Sleep(1 * time.Second)
}
