(function(){
    layui.use(["func", "kline", "wss"], function(){
        var $ = layui.$;
        var symbols = [$(".symbol").attr("data"), $(".stand_symbol").attr("data")];
        
        layui.kline.init("kchart");
        layui.wss.init(symbol);
        layui.func.user_query();
        layui.func.load_assets(symbols, null);
        layui.func.load_open_order(symbol, function(d){
            if(d.ok){
                var tpl = $("#myopenorder-tpl").html();
                
                var rows = d.data;
                for(var i=0; i<rows.length; i++) {
                    layui.laytpl(tpl).render(rows[i], function(html){
                        $(".my-open-order").after(html);
                    })
                }
            }
        });
        layui.func.load_trade_record(symbol, 10, null);

        $(".header .login a").on("click", function(){
            layui.func.login(function(d){
                if(d.ok){
                    layui.func.user_query();
                    //登录完成后加载资产
                    layui.func.load_assets(symbols, null);
                }
            });
        });

        $(".header .logout").on("click", function(){
            layui.func.logout();
            layer.msg("logout sucess");
        });

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
                        layer.msg("order_id: " + d.data.order_id);
                    }else{
                        layer.msg(d.reason);
                    }
                }
            });
        });


        $("body").on("click", ".cancel", function () {
            var me = $(this);
            var order_id = me.parents("tr").attr("order-id");
            layui.func.cancel_order(symbol, order_id, function(d){
                me.parents("tr").remove();
            })
        });


    });





    // layui.use(['laydate', 'layer', 'table', 'element', 'laytpl', 'form', "global"], function () {
    //     var laydate = layui.laydate //日期
    //         , layer = layui.layer //弹层
    //         , table = layui.table //表格
    //         , $ = layui.$
    //         , laytpl = layui.laytpl
    //         , form = layui.form
    //         , element = layui.element; //元素操作 等等...


        
    

    



    //     $("body").on("click", ".cancel", function () {
    //         var me = $(this);
    //         $.ajax({
    //             url: "/api/v1/order/cancel",
    //             type: "post",
    //             dataType: "json",
    //             contentType: "application/json",
    //             data: JSON.stringify({
    //                 "symbol": symbol,
    //                 order_id: me.parents("tr").attr("order-id")
    //             }),
    //             success: function (d) {
    //                 layer.msg("取消 " + d.ok);
    //                 if (d.ok) {
    //                     me.parents("tr").remove();
    //                 }
    //             }
    //         });
    //     });


    //     form.on('select(order_type)', function (data) {
    //         var tt = $(data.elem).attr("data");
    //         if (data.value == "limit") {
    //             $(".item-price-"+tt).show();
    //             $(".item-quantity-"+tt).show();
    //             $(".item-amount-"+tt).hide();
    //             $(".item-market-type-"+tt).hide();
    //             $(".qty-tips-"+tt).hide();
    //         } else if (data.value == "market") {
    //             $(".item-price-"+tt).hide();
    //             $(".item-market-type-"+tt).show();
    //             $(".qty-tips-"+tt).show();
    //         }
    //         form.render('select');
    //     });
    //     form.on('radio(market-type)', function (data) {
    //         if (data.value == "q") {
    //             $(".item-quantity").show();
    //             $(".item-amount").hide();
    //             $(".qty-tips").show();
    //         } else {
    //             $(".item-quantity").hide();
    //             $(".qty-tips").hide();
    //             $(".item-amount").show();
    //         }
    //     });

    //     var load_trade_log = function(){
    //         $.get("/api/v1/trade/record?symbol="+symbol, function (d) {
    //             if (d.ok) {
    //                 var latest_price = d.data[0].price;
    //                 $(".latest-price").html(latest_price);

    //                 var recent_log = d.data.reverse();
    //                 for(var i=0; i<recent_log.length; i++){
    //                     rendertradelog(recent_log[i]);
    //                 }
    //             }
    //         });
    //     };

    //     var load_assets = function(){
    //         var symbols = [$(".symbol").attr("data"), $(".stand_symbol").attr("data")];
    //         $.get("/api/v1/assets/query?symbols=" + symbols.join(","), function(d){
    //             if(d.ok) {
    //                 for(var k in d.data){
    //                     $(".symbol_balance_"+ k).html(d.data[k].available);
    //                 }
    //             }else{
    //                 layer.msg(d.reason);
    //             }
    //         });
    //     };

    //     $().ready(function(){
    //         load_trade_log();
    //         load_assets();
    //     });


        


        

    // });
})()