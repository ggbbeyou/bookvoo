(function(){
    layui.use(['laydate', 'layer', 'table', 'element', 'laytpl', 'form'], function () {
        var laydate = layui.laydate //日期
            , layer = layui.layer //弹层
            , table = layui.table //表格
            , $ = layui.$
            , laytpl = layui.laytpl
            , form = layui.form
            , element = layui.element; //元素操作 等等...


        var socket = function () {
            if (window["WebSocket"]) {
                var protocol = window.location.protocol == "https:" ? "wss:" : "ws:";
                conn = new WebSocket(protocol + "//" + document.location.host + "/ws");
                conn.onclose = function (evt) {
                    layer.msg("<b>WebSocket Connection closed</b>");
                    setTimeout(function () {
                        socket();
                    }, 5e3);
                };

                conn.onmessage = function (evt) {
                    var messages = evt.data.split('\n');
                    for (var i = 0; i < messages.length; i++) {
                        var data = JSON.parse(messages[i]);
                        if (data.tag == "depth") {
                            var info = data.data;
                            var askTpl = $("#depth-ask-tpl").html()
                                , askView = $(".depth-ask")
                                , bidTpl = $("#depth-bid-tpl").html()
                                , bidView = $(".depth-bid");


                            laytpl(askTpl).render(info.ask.reverse(), function (html) {
                                askView.html(html);
                            });
                            laytpl(bidTpl).render(info.bid, function (html) {
                                bidView.html(html);
                            });

                        } else if (data.tag == "trade") {
                            rendertradelog(data.data);
                            
                        } else if (data.tag == "new_order") {
                            var myorderView = $(".myorder"),
                                myorderTpl = $("#myorder-tpl").html();

                            data.data['create_time'] = formatTime(data.data.create_time);
                            laytpl(myorderTpl).render(data.data, function (html) {
                                if ($(".order-item").length > 30) {
                                    $(".order-item").last().remove();
                                }
                                myorderView.after(html);
                            });
                        } else if (data.tag == "latest_price") {
                            $(".latest-price").html(data.data.latest_price);
                        }
                    }
                };

                //订阅一些推送消息
                // conn.send("ping");

            } else {
                layer.msg("<b>Your browser does not support WebSockets.</b>");
            }
        };
        socket();
    })
})()