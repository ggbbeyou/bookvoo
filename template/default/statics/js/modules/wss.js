layui.define(['layer', 'laytpl'], function(exports){
    var conn, 
        $=layui.$,  
        laytpl = layui.laytpl;

    
    layui.config({
        open: '{%',
        close: '%}'
    });

    exports("wss", {
        symbol : "",

        init: function(symbol){
            if(symbol== undefined) {
                layer.msg("wss symbol undefined");
                return;
            }
            this.symbol = symbol;

            this.create();
            this.onopen();
            this.onmessage();
        },

        create: function() {
            if (window["WebSocket"]) {
                var protocol = window.location.protocol == "https:" ? "wss:" : "ws:";
                conn = new WebSocket(protocol + "//" + document.location.host + "/ws");
                conn.onclose = function (evt) {
                    layer.msg("<b>WebSocket Connection closed</b>");
                    setTimeout(function () {
                        socket();
                    }, 5e3);
                };
            } else {
                layer.msg("<b>Your browser does not support WebSockets.</b>");
            }
        },

        onopen: function(){
            //订阅一些推送消息
            var me = this;
            conn.onopen = function(){
                var subs = {
                    sub: [
                        "depth."+ me.symbol,
                        "kline.m1."+me.symbol,
                        "trade.record."+me.symbol,
                    ]
                };
                conn.send(JSON.stringify(subs));
            };
        },

        onmessage: function(){
            var me = this;
            conn.onmessage = function (evt) {
                var messages = evt.data.split('\n');

                console.log(messages);

                for (var i = 0; i < messages.length; i++) {
                    var msg = JSON.parse(messages[i]);
                    if (msg.type == "depth."+me.symbol) {
                        var info = msg.body;
                        var askTpl = $("#depth-ask-tpl").html()
                            , askView = $(".depth-ask")
                            , bidTpl = $("#depth-bid-tpl").html()
                            , bidView = $(".depth-bid");

                        console.log(info);

                        laytpl(askTpl).render(info.ask.reverse(), function (html) {
                            console.log(html);
                            askView.html(html);
                        });
                        laytpl(bidTpl).render(info.bid, function (html) {
                            bidView.html(html);
                        });

                    } else if (msg.type == "trade.record."+ me.symbol) {
                        $(".latest-price").html(msg.body.price);
                        // rendertradelog(msg.body);
                        
                    } else if (msg.type == "new_order."+ me.symbol) {
                        var myorderView = $(".myorder"),
                            myorderTpl = $("#myorder-tpl").html();

                        msg.body['create_time'] = formatTime(msg.body.create_time);
                        laytpl(myorderTpl).render(msg.body, function (html) {
                            if ($(".order-item").length > 30) {
                                $(".order-item").last().remove();
                            }
                            myorderView.after(html);
                        });
                    } else if (msg.type == "latest_price."+ me.symbol) {
                        $(".latest-price").html(msg.body.latest_price);
                    }else if(msg.type=="kline.m1."+ me.symbol){
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
        }



    });
});