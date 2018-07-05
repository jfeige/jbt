# ltools
一个工具类，里面有一些常用的方法



### 方法列表
```
  InArray(obj interface{},target interface{})bool   //检测 指定的value值是否存在Array，Slice，或者Map的value中
  ToMd5(obj interface{})string               
  RandStr(length int,format ...string)string        //随机数(length:字符长度;format:all,char,number)
  Rand(src int)int                                  //返回指定范围内的随机数
  IntToByte(src int)[]byte                          //整型转字节
  ByteToInt(src []byte)int                          //字节转整型
  Reversal(src string)string                        //反转字符串
  ...
  
```
