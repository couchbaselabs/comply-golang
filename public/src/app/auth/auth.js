var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") return Reflect.decorate(decorators, target, key, desc);
    switch (arguments.length) {
        case 2: return decorators.reduceRight(function(o, d) { return (d && d(o)) || o; }, target);
        case 3: return decorators.reduceRight(function(o, d) { return (d && d(target, key)), void 0; }, void 0);
        case 4: return decorators.reduceRight(function(o, d) { return (d && d(target, key, o)) || o; }, desc);
    }
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var core_1 = require("angular2/core");
var http_1 = require("angular2/http");
var router_1 = require("angular2/router");
var authmanager_1 = require("../authmanager");
var utility_1 = require("../utility");
var AuthPage = (function () {
    function AuthPage(http, router, authManager, utility) {
        var _this = this;
        this.router = router;
        this.authManager = authManager;
        this.http = http;
        this.utility = utility;
        this.companies = [];
        this.userCompany = "";
        this.utility.makeGetRequest("/api/company/getAll", []).then(function (result) {
            _this.companies = result;
        }, function (error) {
            console.error(error);
        });
    }
    AuthPage.prototype.changeCompany = function (companyId) {
        console.log(companyId);
    };
    AuthPage.prototype.login = function (email, password) {
        var _this = this;
        if (!email || email == "") {
            console.error("Email must exist");
        }
        else if (!password || password == "") {
            console.error("Password must exist");
        }
        else {
            this.authManager.login(email, password).then(function (result) {
                _this.router.navigate(["Projects"]);
            }, function (error) {
                console.error(error);
            });
        }
    };
    AuthPage.prototype.register = function (firstname, lastname, street, city, state, zip, country, phone, email, password, company) {
        var _this = this;
        var postBody = {
            name: {
                first: firstname,
                last: lastname
            },
            address: {
                street: street,
                city: city,
                state: state,
                zip: zip,
                country: country
            },
            email: email,
            phone: phone,
            password: password,
            company: company
        };
        this.authManager.register(postBody).then(function (result) {
            _this.authManager.login(email, password).then(function (result) {
                _this.router.navigate(["Projects"]);
            }, function (error) {
                console.error(error);
            });
        }, function (error) {
            console.error(error);
        });
        ;
    };
    AuthPage = __decorate([
        core_1.Component({
            selector: "auth",
            viewProviders: [http_1.HTTP_PROVIDERS, authmanager_1.AuthManager, utility_1.Utility]
        }),
        core_1.View({
            templateUrl: "app/auth/auth.html"
        }), 
        __metadata('design:paramtypes', [http_1.Http, router_1.Router, authmanager_1.AuthManager, utility_1.Utility])
    ], AuthPage);
    return AuthPage;
})();
exports.AuthPage = AuthPage;
