package mapreduce

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
)

// doMap does the job of a map worker: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
func doMap(
	jobName string, // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	mapF func(file string, contents string) []KeyValue,
) {
	// TODO:
	// You will need to write this function.
	// You can find the filename for this map task's input to reduce task number
	// r using reduceName(jobName, mapTaskNumber, r). The ihash function (given
	// below doMap) should be used to decide which file a given key belongs into.
	//
	// The intermediate output of a map task is stored in the file
	// system as multiple files whose name indicates which map task produced
	// them, as well as which reduce task they are for. Coming up with a
	// scheme for how to store the key/value pairs on disk can be tricky,
	// especially when taking into account that both keys and values could
	// contain newlines, quotes, and any other character you can think of.
	//
	// One format often used for serializing data to a byte stream that the
	// other end can correctly reconstruct is JSON. You are not required to
	// use JSON, but as the output of the reduce tasks *must* be JSON,
	// familiarizing yourself with it here may prove useful. You can write
	// out a data structure as a JSON string to a file using the commented
	// code below. The corresponding decoding functions can be found in
	// common_reduce.go.
	//
	//   enc := json.NewEncoder(file)
	//   for _, kv := ... {
	//     err := enc.Encode(&kv)
	//
	// Remember to close the file after you have written all the values!
	/*
		在了解了整个流程后，我们来看一下如何实现doMap函数和doReduce函数。在common_map.go文件中有关于doMap函数功能的描述注释，主要操作是打开文件名为inFile的输入文件，读取文件内容，然后调用mapF函数来处理内容，返回值为KeyVaule结构体[common.go]实例，然后生成nReduce个中间文件，提示使用json格式写入。
	*/

	file, _ := ioutil.ReadFile(inFile)
	contents := string(file)
	pairs := mapF(inFile, contents)
	output_maps := make(map[string][]KeyValue)
	for _, kv := range pairs {
		taskNumber := int(ihash(kv.Key)) % nReduce
		outputName := reduceName(jobName, mapTaskNumber, taskNumber) // for each jobtask generate a outputName using reduceName function
		output_maps[outputName] = append(output_maps[outputName], kv)
	}

	// Then we need to create files using the range through theoutput_maps
	for fileName, kvs := range output_maps {
		outputFile, err := os.Open(fileName)
		if os.IsNotExist(err) {
			outputFile, _ = os.Create(fileName)
		}
		enc := json.NewEncoder(outputFile) //NewEncoder returns a new encoder that writes to w.
		err = enc.Encode(&kvs)
		if err != nil {
			fmt.Println("The mapping process is encountering errors:", err)
		}
		outputFile.Close()

	}
}

func ihash(s string) uint32 {
	h := fnv.New32a() // New32a returns a new 32-bit FNV-1a hash.Hash. Its Sum method will lay the value out in big-endian byte order.
	h.Write([]byte(s))
	return h.Sum32()
}
