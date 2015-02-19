var casper = require('casper').create();
var system = require('system');

casper.options.onResourceRequested = function(C, requestData, request) {
    var url = requestData['url'];
    if ((/http:.+?.(gif|png|jpg|woff|ttf|css)/gi).test(url)
        || url.indexOf('googlesyndication') >= 0
        || url.indexOf('googleads.g.doubleclick.net') >= 0
        || url.indexOf('static.googleadsserving.cn') >= 0
        || url.indexOf('cnzz.com') >= 0
        || url.indexOf('wpa.qq.com/') >= 0) {
        request.abort();
    } else {
        //console.log(url);
    }
};

casper.options.onResourceReceived = function(C, response) {

};


casper.userAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 8_0 like Mac OS X) AppleWebKit/600.1.3 (KHTML, like Gecko) Version/8.0 Mobile/12A4345d Safari/600.1.4");

casper.start("http://pachong.org/high.html", function(){

})
.waitUntilVisible(".mainWap", function(){
    items = this.evaluate(function(){
        var trs = document.querySelectorAll(".tb tbody tr");
        ret = [];
        for(var i = 0; i < trs.length; i++){
            var tds = trs[i].querySelectorAll("td");
            ret.push("http://" + tds[1].innerText + ":" + tds[2].innerText);
        }
        return ret;
    });
    for(var i = 0; i < items.length; i++){
        console.log("http://54.223.171.0:7183/register?proxy=" + items[i]);
    }
});

casper.run();
