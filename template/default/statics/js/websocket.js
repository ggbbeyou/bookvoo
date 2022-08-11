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

                 
                

                

            } else {
                layer.msg("<b>Your browser does not support WebSockets.</b>");
            }
        };
        socket();
    })
})()