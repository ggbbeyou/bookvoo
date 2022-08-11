layui.define('layer', function(exports){
    var $ = layui.$;
    exports("kline", {
        init: function(id){
            const chart = klinecharts.init(id);
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
        }
    });
});