layui.define('layer', function(exports){
    var $ = layui.$;
    exports("func", {
        formatTs2Time: function(t) {
            var d = new Date(t*1000);
            return d.getFullYear() + '-' + (d.getMonth() + 1) + '-' + d.getDate() + ' ' + d.getHours() + ':' + d.getMinutes() + ':' + d.getSeconds();
        },


        login: function(callback){
            $.get("/api/v1/login", callback);
        },
        
        load_assets: function(symbols, callback){
            if (symbols.length == 0){
                layer.msg("symbols is length 0");
                return
            }
            $.get("/api/v1/assets/query?symbols=" + symbols.join(","), callback);
        }




    });
});

