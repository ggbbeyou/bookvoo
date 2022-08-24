layui.define('layer', function(exports){
    var $ = layui.$;
    exports("func", {
        formatTs2Time: function(t) {
            var d = new Date(t*1000);
            return d.getFullYear() + '-' + (d.getMonth() + 1) + '-' + d.getDate() + ' ' + d.getHours() + ':' + d.getMinutes() + ':' + d.getSeconds();
        },

        login: function(callback){
            $.get("/api/v1/user/login", callback);
        },

        logout: function(callback){
            $.get("/api/v1/user/logout", function(d){
                window.location.reload();
                if(callback){
                    callback(d);
                }
            });
        },

        cancel_order: function(symbol, order_id, callback){
            $.ajax({
                url: "/api/v1/order/cancel",
                type: "post",
                dataType: "json",
                contentType: "application/json",
                data: JSON.stringify({
                    "symbol": symbol,
                    "order_id": order_id
                }),
                success: callback()
            });
        },
        
        load_assets: function(symbols, callback){
            if (symbols.length == 0){
                layer.msg("symbols is length 0");
                return
            }
            $.get("/api/v1/assets/query?symbols=" + symbols.join(","), function(d){
                if(d.ok) {
                    for(var k in d.data){
                        $(".symbol_balance_"+ k).html(d.data[k].available);
                    }
                }
                
                if(callback){
                    callback(d);
                }
            });
        },


        load_open_order: function(symbol, callback) {
            $.get("/api/v1/order/open?symbol=" + symbol, function(d){
                if(d.ok){}
                if(callback){
                    callback(d);
                }
            });
        },

        load_trade_record: function(symbol, limit, callback){
            var me = this;
            $.get("/api/v1/trade/record?symbol="+symbol+"&limit="+limit, function(d){
                if(d.ok) {
                    var rows = d.data.reverse()
                    for(var i=0; i<rows.length; i++) {
                        me.render_trade_record(rows[i]);
                    }
                }

                if(callback) {
                    callback(d);
                }
            });
        },

        load_latest_price: function(symbol, callback){
            var me = this;
        },

        render_latest_price: function(price) {
            $(".latest-price").html(price);
        },

        render_trade_record: function(data) {
            var logView = $(".trade-log .log"),
                logTpl = $("#trade-record-tpl").html();
        
            data['trade_at'] = this.formatTs2Time(data.trade_at);
            layui.laytpl(logTpl).render(data, function (html) {
                if ($(".log-item").length > 10) {
                    $(".log-item").last().remove();
                }
                logView.after(html);
            });
        },

        render_open_order: function(data){
            var tpl = $("#myopenorder-tpl").html();
            data["create_time"] = this.formatTs2Time(data["create_time"]/1e9);
                    
            layui.laytpl(tpl).render(data, function(html){
                $(".my-open-order").after(html);
            })
        },

        user_query: function(callback){
            $.get("/api/v1/user/query", function(d){
                if(!d.ok){
                    $(".header .login").show();
                    $(".header .userinfo").hide();
                }else{
                    $(".header .login").hide();
                    $(".header .userinfo").show();
                    $(".userinfo .username").html(d.data.username);
                }
                
                if(callback){
                    callback(d);
                }
            });
        }




    });
});

