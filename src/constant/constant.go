package constant

const (
	Output_debug_information bool   = true                                                                                                                                                                                                                              //输出debug信息
	Local_debug              bool   = false                                                                                                                                                                                                                             //进行本地debug，在服务器上使用的事内部网址，在获取devices等信息时不用登陆，但是在本地进行debug的时候使用的是外部网址，需要登陆，所以需要和下面这个cookies配套使用
	Cookie                   string = "eyJhbGciOiJIUzUxMiJ9.eyJhdXQiOlsiVVNFUiJdLCJleHAiOjE1NjU4NjUxODUsImlhdCI6MTU2NTg2MzM4NSwicm9sIjpbIkZhbWlseUN1c3RvbWVyIl0sImp0aSI6IjE5In0.oZFiZR0Bvk5W4dAtEoXoGO5FvQFzhjpTTnRy3o6U0ieYIxLntIi6ayHnlgWVZ0NRZ8kF4GI-HQbjSy2xbDmpBg" //通过在Chrome登陆过后查看cookie并填到这里
)
