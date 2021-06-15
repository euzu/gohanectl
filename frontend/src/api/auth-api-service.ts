import {Observable} from "rxjs";
import ApiService, {DefaultApiService} from "./api-service";
import {first} from "rxjs/operators";

const AUTH_API_URL = 'auth/login';
const AUTH_AUTHENTICATED_API_URL = 'auth/authenticated';

export default interface AuthApiService extends ApiService {
    authenticate(username: string, password: string): Observable<boolean>;
    authenticated(): Observable<boolean>;
}

export class DefaultAuthApiService extends DefaultApiService implements AuthApiService {

    authenticate(username: string, password: string): Observable<boolean> {
        return new Observable<boolean>((observer) => {
            this.post<boolean>(AUTH_API_URL, {username, password}).pipe(first()).subscribe({
                next: (data: any) => { this.setAuthenticated(true, data?.token); observer.next(true)},
                error: (error: any) => {this.setAuthenticated(false, null); observer.next(false)},
                complete: () => {observer.complete()},
            });
        })
    }

    authenticated(): Observable<boolean> {
        return new Observable<boolean>((observable) => {
            this.get<boolean>(AUTH_AUTHENTICATED_API_URL).pipe(first()).subscribe({
                next: (value: boolean) => {
                    DefaultAuthApiService.authenticated$.next(value);
                    observable.next(value)
                },
                error: (err) => observable.error(err),
                complete: () => observable.complete(),
            })
        });
    }
}
