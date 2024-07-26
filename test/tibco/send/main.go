package main

import (
	"github.com/dualm/mq/tibco"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	opt := tibco.TibOption{
		FieldName: "Message",
		Service:   "12120",
		Network:   ";225.12.12.12",
		//Daemon:            []string{},
		TargetSubjectName: "ZXY.FAB.RMS.DEV.EAPsvr",
		//SourceSubjectName: "",
		//PooledMessage:     true,
	}

	infoC := make(chan string)
	errC := make(chan error)

	tib := tibco.NewTibSender(&opt, infoC, errC)

	if err := tib.Run(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := tib.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	//for j := 0; j < 10; j++ {
	//	_n := time.Now()
	//	for i := 0; i < 10; i++ {
	//		err := tib.SendReports([]string{"1", "2", "3", "4", "5"})
	//		if err != nil {
	//			panic(err)
	//		}
	//	}
	//	log.Println(time.Since(_n).Nanoseconds())
	//}

	re1, err1 := tib.SendRequest(str, 10000)
	if err1 != nil {
		log.Fatal(err1)
	}

	re2, err1 := re1.GetString("Message", 0)
	if err1 != nil {
		log.Fatal(err1)
	}

	log.Println(re2)

	os.Exit(0)

	k := 0
	for {
		n := rand.Intn(5)
		i := 0
		for {
			if i >= n {
				break
			}
			i++
			time.Sleep(time.Second)
		}
		k += 1
		log.Println(k)
		re, err := tib.SendRequest(`
<?xml version="1.0" encoding="UTF-8"?>
<Message>
    <Header>
        <MESSAGENAME>ALLRecipeParameterUploadRequest</MESSAGENAME>
        <TRANSACTIONID>20231026201805509000</TRANSACTIONID>
        <ORIGINALSOURCESUBJECTNAME />
        <SOURCESUBJECTNAME>ZXY.FAB.RMS.DEV.RMSsvr</SOURCESUBJECTNAME>
        <TARGETSUBJECTNAME>ZXY.FAB.RMS.DEV.EAPsvr</TARGETSUBJECTNAME>
        <FACTORY />
        <AREA />
        <MACHINENAME />
        <UNITNAME />
        <SUBUNITNAME />
        <EVENTUSER>EPG01</EVENTUSER>
        <EVENTCOMMENT>RecipeUploadRequest</EVENTCOMMENT>
    </Header>
    <Body>
        <RECIPENAME>TEST01</RECIPENAME>
        <MACHINENAME>EPG01</MACHINENAME>
        <FLAG />
        <RECIPETYPE />
    </Body>
</Message>
`, 100000000)
		if err != nil {
			log.Fatal(err)
		}

		_, err = re.GetString("Message", 0)
		if err != nil {
			log.Fatal(err)
		}
	}
}

var str = `
<Message>
    <Header>
        <MESSAGENAME>RecipeDownload</MESSAGENAME>
        <TRANSACTIONID>20231026202307993000</TRANSACTIONID>
        <ORIGINALSOURCESUBJECTNAME />
        <SOURCESUBJECTNAME>ZXY.FAB.RMS.DEV.RMSsvr</SOURCESUBJECTNAME>
        <TARGETSUBJECTNAME>ZXY.FAB.RMS.DEV.EAPsvr</TARGETSUBJECTNAME>
        <FACTORY />
        <AREA />
        <MACHINENAME />
        <UNITNAME />
        <SUBUNITNAME />
        <EVENTUSER>admin</EVENTUSER>
        <EVENTCOMMENT>3</EVENTCOMMENT>
    </Header>
    <Body>
        <UNITNAME />
        <SUBUNITNAME />
        <RECIPENAME>TEST001</RECIPENAME>
        <MACHINENAME>EPG01</MACHINENAME>
        <RECIPETYPE>Process</RECIPETYPE>
        <STEPLIST>
            <STEP>
                <STEPNAME>Start LTP growth</STEPNAME>
                <STEPRECIPEBODY>[LPT_Time]	"Start LTP growth",&#xD;
		begin stat stat_all,&#xD;
		StepCode = 150,&#xD;
                TMGa_1.run = open,&#xD;
		TMAl_1.run = close,&#xD;
		Cp2Mg_1.run = open,&#xD;
		Cp2Mg_2.run = open;</STEPRECIPEBODY>
                <PARAMETERLIST>
                    <PARAMETER>
                        <ITEMTYPE>null</ITEMTYPE>
                        <ITEMVALUE>07:40</ITEMVALUE>
                        <ITEMNAME>LPT_Time</ITEMNAME>
                    </PARAMETER>
                </PARAMETERLIST>
                <STEPNO>64</STEPNO>
            </STEP>
        </STEPLIST>
    </Body>
</Message>`

var str2 = `
<Message>
    <Header>
        <MESSAGENAME>RecipeDownload</MESSAGENAME>
        <TRANSACTIONID>20231026202307993000</TRANSACTIONID>
        <ORIGINALSOURCESUBJECTNAME />
        <SOURCESUBJECTNAME>ZXY.FAB.RMS.DEV.RMSsvr</SOURCESUBJECTNAME>
        <TARGETSUBJECTNAME>ZXY.FAB.RMS.DEV.EAPsvr</TARGETSUBJECTNAME>
        <FACTORY />
        <AREA />
        <MACHINENAME />
        <UNITNAME />
        <SUBUNITNAME />
        <EVENTUSER>admin</EVENTUSER>
        <EVENTCOMMENT>3</EVENTCOMMENT>
    </Header>
    <Body>
        <UNITNAME />
        <SUBUNITNAME />
        <RECIPENAME>TEST001</RECIPENAME>
        <MACHINENAME>EPG01</MACHINENAME>
        <STEPLIST />
        <RECIPETYPE>Process</RECIPETYPE>
        <RECIPELIST />
    </Body>
</Message>
`
