package api_caller

import (
	"fmt"
	"quicky-go/enpoint_controller"
	"reflect"
	"runtime"
	"strings"
)

// Định nghĩa struct Invoker
type Invoker struct {
	EndPoint string                  `json:"enpoint"` // Sử dụng tag json để map trường JSON
	Params   *map[string]interface{} `json:"params"`  // nhận vào một map chứa các tham số cần gửi lên API ví dụ {"name": "John", "age": 30}
}

func (i *Invoker) Invoke() (interface{}, error) {
	apis := strings.Split(i.EndPoint, ".")
	packName := apis[0]
	funcName := apis[1]
	ret, ex := CallFuncByName(packName, funcName)
	if ex != nil {
		return nil, ex
	}
	return ret, nil
}

// GọiFuncByName gọi một hàm trong một gói bằng tên gói và tên hàm.
// Hàm này giả định rằng hàm không có tham số và trả về một giá trị.
func CallFuncByName(packageName, funcName string) (interface{}, error) {
	// 1. Lấy con trỏ đến hàm.
	// Dựa vào quy ước đặt tên của Go, tên hàm có thể khác với tên được export.
	// Ví dụ: hàm "myFunc" có thể được export là "MyFunc".
	// Ta cần tìm đúng tên hàm đã được export.
	fn := reflect.ValueOf(GetFunctionPointer(packageName, funcName))
	if !fn.IsValid() {
		return nil, fmt.Errorf("function %s in package %s not found", funcName, packageName)
	}

	// 2. Kiểm tra loại của giá trị trả về.
	fnType := fn.Type()
	if fnType.NumIn() != 0 {
		return nil, fmt.Errorf("function %s in package %s expects %d arguments, but 0 were provided", funcName, packageName, fnType.NumIn())
	}

	// 3. Gọi hàm.
	// Tạo slice rỗng cho các argument vì hàm không có tham số.
	args := []reflect.Value{}
	results := fn.Call(args)

	// 4. Xử lý kết quả trả về.
	if len(results) == 0 {
		return nil, nil // Hàm không trả về gì.
	}
	if len(results) == 1 {
		return results[0].Interface(), nil // Hàm trả về một giá trị.
	}

	//Hàm trả về nhiều hơn 1 giá trị
	vals := make([]interface{}, len(results))
	for i, v := range results {
		vals[i] = v.Interface()
	}
	return vals, nil
	//return nil, fmt.Errorf("function %s in package %s returns multiple values", funcName, packageName) // Hàm trả về nhiều giá trị.
}

// GetFunctionPointer lấy con trỏ hàm từ tên gói và tên hàm.
func GetFunctionPointer(packageName, funcName string) interface{} {
	// reflect.ValueOf trả về một Value đại diện cho giá trị của một biến.
	// Hàm này có thể là một biến, một hằng số, một hàm, v.v.
	// Trong trường hợp này, chúng ta muốn lấy Value của một hàm.

	// runtime.FuncForPC trả về một Func đại diện cho hàm bắt đầu tại địa chỉ chương trình pc.
	// reflect.ValueOf(fn).Pointer() trả về địa chỉ chương trình của hàm fn.
	// Caller(0) trả về program counter (pc) của hàm gọi hiện tại.

	pc := reflect.ValueOf(enpoint_controller.GetEnpointList).Pointer()
	fn := runtime.FuncForPC(pc)

	if fn == nil {
		return nil
	}
	// Lấy tên file và số dòng của hàm
	fileName, _ := fn.FileLine(pc)

	//packageName from file name
	parts := strings.Split(fileName, "/")
	if len(parts) < 2 {
		return nil // Không tìm thấy dấu gạch chéo
	}
	endPointPackageName := parts[len(parts)-2]
	packageName = parts[len(parts)-2]

	targetFuncName := funcName
	if !strings.HasPrefix(targetFuncName, strings.ToLower(endPointPackageName)) {
		targetFuncName = strings.ToLower(packageName) + "." + funcName
	}
	fmt.Println("targetFuncName", targetFuncName)
	// Hàm runtime.FuncForPC trả về một *runtime.Func,
	//  đại diện cho thông tin về một hàm đang thực thi.
	//  Phương thức Name() của *runtime.Func trả về tên của hàm đó.
	//  Tên này bao gồm tên gói.
	// Ví dụ:
	//  - Nếu bạn gọi một hàm trong main package, nó có thể trả về "main.main".
	//  - Nếu bạn gọi một hàm trong package "mypackage", nó có thể trả về "mypackage.MyFunc".

	// Lặp qua tất cả các hàm trong chương trình.
	for i := 0; i < 10000; i++ { // Giới hạn để tránh lặp vô hạn.
		fn := runtime.FuncForPC(pc + uintptr(i))
		if fn == nil {
			break // Không còn hàm nào để kiểm tra.
		}
		name := fn.Name()
		//fmt.Println("func name",name)
		// Kiểm tra xem tên hàm có khớp với tên hàm chúng ta cần tìm không.
		//if name == fmt.Sprintf("%s.%s", packageName, funcName) {
		if name == targetFuncName {
			// Trả về Value của hàm.
			return reflect.ValueOf(fn).Interface()
		}
	}
	return nil // Không tìm thấy hàm.
}
