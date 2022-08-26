/*
 * @Date: 2022-08-23 17:02:41
 */
/**
 * [获取URL中的参数名及参数值的集合]
 * 示例URL:http://htmlJsTest/getrequest.html?uid=admin&rid=1&fid=2&name=小明
 * @param {[string]} urlStr [当该参数不为空的时候，则解析该url中的参数集合]
 * @return {[string]}       [参数集合]
 */
function GetRequest(urlStr) {
    if (typeof urlStr == "undefined") {
        var url = decodeURI(location.search); //获取url中"?"符后的字符串
    } else {
        var url = (typeof urlStr.split("?")[1] == 'undefined' ? "" : "?" + urlStr.split("?")[1]);
    }
    var theRequest = new Object();
    if (url.indexOf("?") != -1) {
        var str = url.substr(1);
        strs = str.split("&");
        for (var i = 0; i < strs.length; i++) {
            theRequest[strs[i].split("=")[0]] = decodeURI(strs[i].split("=")[1]);
        }
    }
    return theRequest;
}

/**
 * 追加参数（如果有，则替换）
 */
function ReplaceOrAddRequest(params, url = "")
{
    if (url == "") {
        url = window.location.href
        baseUrl = window.location.origin + window.location.pathname + "?";
    } else {
        baseUrl = url;
    }
    // 获取当前的参数
    currentRequest = GetRequest(url);

    $.each(params, function(_, item) {
        currentRequest[item.name] = item.value
    });

    newRequest = new URLSearchParams(currentRequest).toString()
    newUrl = baseUrl + newRequest
    return newUrl;
}

function GetUrlParam(name)
{
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
    var r = window.location.search.substr(1).match(reg);  //匹配目标参数
    if (r != null) return unescape(r[2]); return null; //返回参数值
}
