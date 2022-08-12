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

