(function(){
    layui.use(['laydate', 'layer', 'table', 'element', 'laytpl', 'form'], function () {
        var laydate = layui.laydate //日期
            , layer = layui.layer //弹层
            , table = layui.table //表格
            , $ = layui.$
            , laytpl = layui.laytpl
            , form = layui.form
            , element = layui.element; //元素操作 等等...

        laytpl.config({
            open: '{%',
            close: '%}'
        });

        
        var kchart = function(){
            const chart = klinecharts.init("kchart");
            chart.createTechnicalIndicator('MA', false, { id: 'candle_pane' });
            // Create sub technical indicator VOL
            chart.createTechnicalIndicator('VOL');
            window.Kdata = [];
            var loadKline = function(){
                $.ajax({
                    url: "/api/v1/market/klines?symbol="+symbol+"&period=m1",
                    type: "GET",
                    dataType: "json",
                    success: function(d){
                        
                        for(var i=0; i<d.length; i++){
                            Kdata.push({
                                timestamp: d[i][0] * 1000,
                                open: parseFloat(d[i][1]),
                                high: parseFloat(d[i][2]),
                                low: parseFloat(d[i][3]),
                                close: parseFloat(d[i][4]),
                                volume: parseFloat(d[i][5]),
                            });
                        }

                        chart.applyNewData(Kdata.reverse());
                    },
                })
            };

            loadKline();
            window.kLchart = chart;
        };
        kchart();

        function rendertradelog(data) {
            var logView = $(".trade-log .log"),
                logTpl = $("#trade-log-tpl").html();
        
            data['TradeTime'] = formatTime(data.TradeTime);
            laytpl(logTpl).render(data, function (html) {
                if ($(".log-item").length > 10) {
                    $(".log-item").last().remove();
                }
                logView.after(html);
        
                //remove myorder
                $("tr[order-id='" + data.AskOrderId + "']").remove();
                $("tr[order-id='" + data.BidOrderId + "']").remove();
            });
        }



        function createUUID() {
            var dt = new Date().getTime();
            var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
                var r = (dt + Math.random() * 16) % 16 | 0;
                dt = Math.floor(dt / 16);
                return (c == 'x' ? r : (r & 0x3 | 0x8)).toString(16);
            });
            return uuid;
        }

        /*layer弹出一个示例*/
        //   layer.msg('Hello World');

        $(".opt").on("click", function () {
            var type = $(this).hasClass("sell") ? "ask" : "bid";
            var price_type = $("select[name='price_type_"+ type +"']").val();
            var mtype = $("input[name='mtype_"+ type +"']:checked").val();

            $.ajax({
                url: "/api/new_order",
                type: "post",
                dataType: "json",
                contentType: "application/json",
                data: function () {
                    var data = {
                        price_type: price_type,
                        order_type: type,
                        "symbol": symbol,
                    };

                    if (price_type == "market") {
                        if (mtype == "q") {
                            data.quantity = $("input[name='quantity_"+ type +"']").val();
                        } else {
                            data.amount = $("input[name='amount_"+ type +"']").val();
                        }
                    } else {
                        data.price = $("input[name='price_"+ type +"']").val();
                        data.quantity = $("input[name='quantity_"+ type +"']").val();
                    }

                    console.log(data);
                    return JSON.stringify(data)
                }(),
                success: function (d) {
                    if(d.ok){
                        layer.msg("下单" + d.ok + " askLen:" + d.data.ask_len + " bidLen:" + d.data.bid_len);
                    }else{
                        layer.msg(d.error);
                    }
                }
            });
        });

        $(".test-rand").on("click", function () {
            var op_type = "ask", me = $(this);
            if ($(this).hasClass("buy")) {
                op_type = "bid";
            }

            me.attr("disabled", true);

            $.ajax({
                url: "/api/test_rand?op_type=" + op_type + "&symbol="+symbol,
                type: "get",
                success: function (d) {
                    layer.msg("操作" + d.ok + " askLen:" + d.data.ask_len + " bidLen:" + d.data.bid_len);
                    me.attr("disabled", false);
                }
            });

        });

        $("body").on("click", ".cancel", function () {
            var me = $(this);
            $.ajax({
                url: "/api/cancel_order",
                type: "post",
                dataType: "json",
                contentType: "application/json",
                data: JSON.stringify({
                    "symbol": symbol,
                    order_id: me.parents("tr").attr("order-id")
                }),
                success: function (d) {
                    layer.msg("取消 " + d.ok);
                    if (d.ok) {
                        me.parents("tr").remove();
                    }
                }
            });
        });


        form.on('select(price_type)', function (data) {
            var tt = $(data.elem).attr("data");
            if (data.value == "limit") {
                $(".item-price-"+tt).show();
                $(".item-quantity-"+tt).show();
                $(".item-amount-"+tt).hide();
                $(".item-market-type-"+tt).hide();
                $(".qty-tips-"+tt).hide();
            } else if (data.value == "market") {
                $(".item-price-"+tt).hide();
                $(".item-market-type-"+tt).show();
                $(".qty-tips-"+tt).show();
            }
            form.render('select');
        });
        form.on('radio(market-type)', function (data) {
            if (data.value == "q") {
                $(".item-quantity").show();
                $(".item-amount").hide();
                $(".qty-tips").show();
            } else {
                $(".item-quantity").hide();
                $(".qty-tips").hide();
                $(".item-amount").show();
            }
        });


        $().ready(function(){
            $.get("/api/trade_log?symbol="+symbol, function (d) {
                if (d.ok) {
                    $(".latest-price").html(d.data.latest_price);

                    var recent_log = d.data.trade_log;
                    for(var i=0; i<recent_log.length; i++){
                        rendertradelog(recent_log[i]);
                    }

                }
            });
        });


        


        

    });
})()