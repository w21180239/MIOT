package constant

const (
	Output_debug_information bool   = true                                                                                                                                                                                                                              //输出debug信息
	Local_debug              bool   = true                                                                                                                                                                                                                              //进行本地debug，在服务器上使用的事内部网址，在获取devices等信息时不用登陆，但是在本地进行debug的时候使用的是外部网址，需要登陆，所以需要和下面这个cookies配套使用
	Cookie                   string = "eyJhbGciOiJIUzUxMiJ9.eyJhdXQiOlsiVVNFUiJdLCJleHAiOjE1NjU4NjAyODAsImlhdCI6MTU2NTg1ODQ4MCwicm9sIjpbIkZhbWlseUN1c3RvbWVyIl0sImp0aSI6IjE5In0.TTMi85AwTpIOljx-44IPMalHY-MunTTwhRAZW1Kdt_90MDxOe8zy9_5n_D1mc6YMhRtHeA-sd8ySrLysm919zg" //通过在Chrome登陆过后查看cookie并填到这里
)
