(function(){
    layui.use(['laydate', 'layer', 'table', 'element', 'laytpl', 'form'], function () {
        var laydate = layui.laydate //日期
            , layer = layui.layer //弹层
            , table = layui.table //表格
            , $ = layui.$
            , laytpl = layui.laytpl
            , form = layui.form
            , element = layui.element; //元素操作 等等...

        var cur_symbol = symbol;

        function rendertradelog(data) {
            var logView = $(".trade-log .log"),
                logTpl = $("#trade-log-tpl").html();
        
            
            data['trade_at'] = formatTs2Time(parseInt(data.trade_at/1e9));
            

            laytpl(logTpl).render(data, function (html) {
                if ($(".log-item").length > 10) {
                    $(".log-item").last().remove();
                }
                logView.after(html);
            });
        }

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
                        var msg = JSON.parse(messages[i]);
                        if (msg.type == "depth."+cur_symbol) {
                            var info = msg.body;
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

                        } else if (msg.type == "trade.record."+cur_symbol) {
                            $(".latest-price").html(msg.body.price);
                            rendertradelog(msg.body);
                            
                        } else if (msg.type == "new_order."+cur_symbol) {
                            var myorderView = $(".myorder"),
                                myorderTpl = $("#myorder-tpl").html();

                            msg.body['create_time'] = formatTime(msg.body.create_time);
                            laytpl(myorderTpl).render(msg.body, function (html) {
                                if ($(".order-item").length > 30) {
                                    $(".order-item").last().remove();
                                }
                                myorderView.after(html);
                            });
                        } else if (msg.type == "latest_price."+cur_symbol) {
                            $(".latest-price").html(msg.body.latest_price);
                        }else if(msg.type=="kline.m1."+cur_symbol){
                            window.kLchart.updateData({
                                timestamp: msg.body.open_at * 1000,
                                open: parseFloat(msg.body.open),
                                high: parseFloat(msg.body.high),
                                low: parseFloat(msg.body.low),
                                close: parseFloat(msg.body.close),
                                volume: parseFloat(msg.body.volume),
                            });
                        }
                    }
                };

            } else {
                layer.msg("<b>Your browser does not support WebSockets.</b>");
            }
        };
        socket();
    })
})()