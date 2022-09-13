package jmeter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func DisableJmeterResultTrees(scriptFile string) { // 读取,处理,写入
	// 逐行读取
	file, err := os.Open(scriptFile)
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(file)
	buf := strings.Builder{}

	for {
		line, _, err := reader.ReadLine()
		str := string(line)
		if err == io.EOF {
			break
		}
		// 逐行替换
		for k, v := range jmeterTreeResultMap() {
			if strings.TrimSpace(k) == strings.TrimSpace(str) {
				log.Println(str)
				log.Println("replace: ", v)
				str = v
			}
		}
		buf.WriteString(str)
	}
	file.Close()
	//写入

	newFile, err := os.Create(scriptFile)
	if err != nil {
		fmt.Println(err)
	}
	_, err = newFile.WriteString(buf.String())
	if err != nil {
		fmt.Println(err)
	}
	newFile.Close()
}

func GetJmxFilename(Cmd string) (string, error) { // 获取jmx filename
	list := strings.Split(Cmd, " ")
	for _, v := range list {
		if strings.HasSuffix(v, ".jmx") {
			return v, nil
		}
	}
	return "", os.ErrExist
}

func jmeterTreeResultMap() map[string]string {
	res := make(map[string]string, 0)
	res["<ResultCollector guiclass=\"ViewResultsFullVisualizer\" testclass=\"ResultCollector\" testname=\"View Results Tree\" enabled=\"true\">"] = "<ResultCollector guiclass=\"ViewResultsFullVisualizer\" testclass=\"ResultCollector\" testname=\"View Results Tree\" enabled=\"false\">"
	res["<ResultCollector guiclass=\"SummaryReport\" testclass=\"ResultCollector\" testname=\"Summary Report\" enabled=\"true\">"] = "<ResultCollector guiclass=\"SummaryReport\" testclass=\"ResultCollector\" testname=\"Summary Report\" enabled=\"false\">"
	res["<ResultCollector guiclass=\"StatVisualizer\" testclass=\"ResultCollector\" testname=\"Aggregate Report\" enabled=\"true\">"] = "<ResultCollector guiclass=\"StatVisualizer\" testclass=\"ResultCollector\" testname=\"Aggregate Report\" enabled=\"false\">"
	res["<BackendListener guiclass=\"BackendListenerGui\" testclass=\"BackendListener\" testname=\"Backend Listener\" enabled=\"true\">"] = "<BackendListener guiclass=\"BackendListenerGui\" testclass=\"BackendListener\" testname=\"Backend Listener\" enabled=\"false\">"
	res["<ResultCollector guiclass=\"StatGraphVisualizer\" testclass=\"ResultCollector\" testname=\"Aggregate Graph\" enabled=\"true\">"] = "<ResultCollector guiclass=\"StatGraphVisualizer\" testclass=\"ResultCollector\" testname=\"Aggregate Graph\" enabled=\"false\">"
	res["<ResultCollector guiclass=\"AssertionVisualizer\" testclass=\"ResultCollector\" testname=\"Assertion Results\" enabled=\"true\">"] = "<ResultCollector guiclass=\"AssertionVisualizer\" testclass=\"ResultCollector\" testname=\"Assertion Results\" enabled=\"false\">"
	res["<ResultCollector guiclass=\"ComparisonVisualizer\" testclass=\"ResultCollector\" testname=\"Comparison Assertion Visualizer\" enabled=\"true\">"] = "<ResultCollector guiclass=\"ComparisonVisualizer\" testclass=\"ResultCollector\" testname=\"Comparison Assertion Visualizer\" enabled=\"false\">"
	res["<Summariser guiclass=\"SummariserGui\" testclass=\"Summariser\" testname=\"Generate Summary Results\" enabled=\"true\"/>"] = "<Summariser guiclass=\"SummariserGui\" testclass=\"Summariser\" testname=\"Generate Summary Results\" enabled=\"false\"/>"
	res["<ResultCollector guiclass=\"GraphVisualizer\" testclass=\"ResultCollector\" testname=\"Graph Results\" enabled=\"true\">"] = "<ResultCollector guiclass=\"GraphVisualizer\" testclass=\"ResultCollector\" testname=\"Graph Results\" enabled=\"false\">"
	res["<JSR223Listener guiclass=\"TestBeanGUI\" testclass=\"JSR223Listener\" testname=\"JSR223 Listener\" enabled=\"true\">"] = "<JSR223Listener guiclass=\"TestBeanGUI\" testclass=\"JSR223Listener\" testname=\"JSR223 Listener\" enabled=\"false\">"
	res["<MailerResultCollector guiclass=\"MailerVisualizer\" testclass=\"MailerResultCollector\" testname=\"Mailer Visualizer\" enabled=\"true\">"] = "<MailerResultCollector guiclass=\"MailerVisualizer\" testclass=\"MailerResultCollector\" testname=\"Mailer Visualizer\" enabled=\"false\">"
	res["<ResultSaver guiclass=\"ResultSaverGui\" testclass=\"ResultSaver\" testname=\"Save Responses to a file\" enabled=\"true\">"] = "<ResultSaver guiclass=\"ResultSaverGui\" testclass=\"ResultSaver\" testname=\"Save Responses to a file\" enabled=\"false\">"
	res["<ResultCollector guiclass=\"RespTimeGraphVisualizer\" testclass=\"ResultCollector\" testname=\"Response Time Graph\" enabled=\"true\">"] = "<ResultCollector guiclass=\"RespTimeGraphVisualizer\" testclass=\"ResultCollector\" testname=\"Response Time Graph\" enabled=\"false\">"
	res["<ResultCollector guiclass=\"SimpleDataWriter\" testclass=\"ResultCollector\" testname=\"Simple Data Writer\" enabled=\"true\">"] = "<ResultCollector guiclass=\"SimpleDataWriter\" testclass=\"ResultCollector\" testname=\"Simple Data Writer\" enabled=\"false\">"
	res["<ResultCollector guiclass=\"TableVisualizer\" testclass=\"ResultCollector\" testname=\"View Results in Table\" enabled=\"true\">"] = "<ResultCollector guiclass=\"TableVisualizer\" testclass=\"ResultCollector\" testname=\"View Results in Table\" enabled=\"false\">"
	res["<BeanShellListener guiclass=\"TestBeanGUI\" testclass=\"BeanShellListener\" testname=\"BeanShell Listener\" enabled=\"true\">"] = "<BeanShellListener guiclass=\"TestBeanGUI\" testclass=\"BeanShellListener\" testname=\"BeanShell Listener\" enabled=\"false\">"
	return res
}
