layui.define('layer', function(exports){
    var conn;
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
        }



    });
});