import AuthApiService, {DefaultAuthApiService} from "../api/auth-api-service";
import {Observable} from "rxjs";

export default class AuthService {

    constructor(private authApiService: AuthApiService = new DefaultAuthApiService()) {
    }

    isAuthenticated(): Observable<boolean> {
        return this.authApiService.isAuthenticated();
    }

    login(username: string, password: string): Observable<boolean> {
        return this.authApiService.authenticate(username, password);
    }

    authenticated(): Observable<boolean> {
        return this.authApiService.authenticated();
    }

}
