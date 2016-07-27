import {Injectable, Inject} from "@angular/core";
import {Http, Request, RequestMethod, Headers, HTTP_PROVIDERS} from "@angular/http";

@Injectable()

export class Utility {

    http: Http;

    constructor(http: Http) {
        this.http = http;
    }

    makePostRequest(url: string, params: Array<string>, body: Object, query?: Object) {
        var fullUrl: string = url;
        if(params && params.length > 0) {
            let encodedParams: Array<string> = [];
            for(let i = 0; i < params.length; i++) {
                encodedParams.push(encodeURIComponent(params[i]));
            }
            fullUrl = fullUrl + "/" + encodedParams.join("/");
        }
        if(query) {
            fullUrl += this.jsonToQueryString(query);
        }
        console.log("DEBUG: POST FULL URL: ", fullUrl, " BODY: ", JSON.stringify(body));
        return new Promise((resolve, reject) => {
            var requestHeaders = new Headers();
            requestHeaders.append("Content-Type", "application/json");
            this.http.request(new Request({
                method: RequestMethod.Post,
                url: fullUrl,
                body: JSON.stringify(body),
                headers: requestHeaders
            }))
            .subscribe((success) => {
                console.log("DEBUG: POST RESPONSE: ", fullUrl, ": ", success.json());
                resolve(success.json());
            }, (error) => {
                reject(error.json());
            });
        });
    }

    makeGetRequest(url: string, params: Array<string>, query?: Object) {
        var fullUrl: string = url;
        if(params && params.length > 0) {
            let encodedParams: Array<string> = [];
            for(let i = 0; i < params.length; i++) {
                encodedParams.push(encodeURIComponent(params[i]));
            }
            fullUrl = fullUrl + "/" + encodedParams.join("/");
        }
        if(query) {
            fullUrl += this.jsonToQueryString(query);
        }
        console.log("DEBUG: GET FULL URL: ", fullUrl);
        return new Promise((resolve, reject) => {
            this.http.get(fullUrl)
            .subscribe((success) => {
                console.log("DEBUG: GET RESPONSE: ", fullUrl, ": ", success.json());
                resolve(success.json());
            }, (error) => {
                reject(error.json());
            });
        });
    }

    private jsonToQueryString(json: Object) {
        return '?' +
            Object.keys(json).map(function(key) {
                return encodeURIComponent(key) + '=' +
                encodeURIComponent(json[key]);
            }).join('&');
    }

}
