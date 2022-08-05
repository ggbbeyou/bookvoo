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
        
            data['trade_at'] = formatTs2Time(data.trade_at);
            laytpl(logTpl).render(data, function (html) {
                if ($(".log-item").length > 10) {
                    $(".log-item").last().remove();
                }
                logView.after(html);
            });
        }


        $(".opt").on("click", function () {
            var side = $(this).hasClass("sell") ? "sell" : "buy";
            var order_type = $("select[name='order_type_"+ side +"']").val();
            var mtype = $("input[name='mtype_"+ side +"']:checked").val();

            $.ajax({
                url: "/api/v1/order/new",
                type: "post",
                dataType: "json",
                contentType: "application/json",
                data: function () {
                    var data = {
                        "order_type": order_type,
                        "side": side,
                        "symbol": symbol,
                    };

                    if (order_type == "market") {
                        if (mtype == "q") {
                            data.quantity = $("input[name='quantity_"+ side +"']").val();
                        } else {
                            data.amount = $("input[name='amount_"+ side +"']").val();
                        }
                    } else {
                        data.price = $("input[name='price_"+ side +"']").val();
                        data.quantity = $("input[name='quantity_"+ side +"']").val();
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
                url: "/api/v1/order/cancel",
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


        form.on('select(order_type)', function (data) {
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
            $.get("/api/v1/trade/record?symbol="+symbol, function (d) {
                if (d.ok) {
                    var recent_log = d.data.reverse();
                    for(var i=0; i<recent_log.length; i++){
                        rendertradelog(recent_log[i]);
                    }
                }
            });
        });


        


        

    });
})()