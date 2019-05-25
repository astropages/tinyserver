/*
	数据包模块测试
*/

package tsnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//TestDataPack 测试datapack模块（封包、打包）
func TestDataPack(t *testing.T) {
	fmt.Println("DataPack模块测试")

	/*
		模拟服务器（解包）
	*/

	listenner, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("net.Listen error", err)
		return
	}
	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("listenner.Accept error", err)
			}

			//读业务
			go func(conn *net.Conn) {
				//先读取客户端数据：拆包过程 |Datalen|ID|data|
				dp := NewDataPack() //创建一个数据包对象
				for {
					//读取head部分
					head := make([]byte, dp.GetHeadLen()) //8个字节
					_, err := io.ReadFull(*conn, head)    //ReadFull读取直到head填满才会返回，否则会阻塞
					if err != nil {
						fmt.Println("io.ReadFull error", err)
						return
					}

					//解包

					//先把数据填充到Datalen和ID属性中
					msg, err := dp.UnPack(head)
					if err != nil {
						fmt.Println("dp.UnPack error ", err)
						return
					}
					//判断是否有数据，如果数据有内容则需要二次读取
					if msg.GetMsgLen() > 0 {
						//由于msg是IMessage类型，因此需要通过断言转换为Message对象
						message := msg.(*Message)
						//给message开辟空间，长度为head中保存的值
						message.Data = make([]byte, message.GetMsgLen())
						//根据数据长度进行读取
						_, err := io.ReadFull(*conn, message.Data)
						if err != nil {
							fmt.Println("io.ReadFull error ", err)
							return
						}

						//打印数据
						fmt.Printf("数据%d: %s (长度为%d)\n", message.ID, string(message.Data), message.Datalen)
					}
				}
			}(&conn)
		}
	}()

	/*
		模拟客户端（封包）
	*/

	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return
	}

	dp := NewDataPack() //创建一个数据包对象

	//封包

	//包1
	msg1 := &Message{
		ID:      1,
		Datalen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	msg1pack, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("dp.Pack error", err)
		return
	}

	//包2
	msg2 := &Message{
		ID:      2,
		Datalen: 4,
		Data:    []byte{'t', 'e', 's', 't'},
	}
	msg2pack, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("dp.Pack error", err)
		return
	}

	//发送
	msg1pack = append(msg1pack, msg2pack...) //将两个数据包粘一起发送测试
	conn.Write(msg1pack)

	//防止测试函数结束
	select {}
}
