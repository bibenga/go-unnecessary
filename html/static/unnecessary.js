async function UnnecessaryInvoke(behaviorId, eventCode, elementId) {
    try {
        var link = new URL(window.location.href);
        link.searchParams.set('anticache', '' + Math.random());
        var response = await fetch(link, {
            method: "POST",
            credentials: "same-origin",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "behaviorId": behaviorId,
                "eventCode": eventCode,
                "elementId": elementId,
            }),
        });
        var ajaxResponse = await response.json();
        ajaxResponse.items.forEach(item => {
            console.log(item);
            el = document.getElementById(item.markupId);
            if (!!item.body) {
                el.outerHTML = item.body;
            }
            if (!!item.script) {
                // todo: insecure
                eval(item.script);
            }
        });
    } catch (e) {
        console.log(e);
    }
}

function UnnecessaryAddEventListener(lazy, behaviorId, eventCode, elementId) {
    var link = () => {
        var el = document.getElementById(elementId);
        el.addEventListener(eventCode, event => {
            UnnecessaryInvoke(behaviorId, eventCode, elementId)
            event.preventDefault();
            event.stopPropagation();
            return false;
        });
    }
    if (lazy) {
        document.addEventListener("DOMContentLoaded", link);
    } else {
        link();
    }
}
