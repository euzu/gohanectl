import config from '../config';
import {BehaviorSubject, Observable} from "rxjs";
import axios from "axios";

const KEY_AUTH_TOKEN = 'auth-token';
const HEADER_AUTHORIZATION = 'Authorization';
const HEADER_CONTENT_TYPE = 'Content-Type';
const HEADER_LANGUAGE = 'X-Language';
const HEADER_ACCEPT = 'Accept';

export default interface ApiService {
    isAuthenticated(): Observable<boolean>;
    setAuthenticated(value: boolean, token: string) : void;
    // get<T>(query: string, url?: string): Observable<T>;
    // post<T>(query: string, payload: any, url?: string): Observable<T>;
    // put<T>(query: string, payload: any, url?: string): Observable<T>;
    // delete<T>(query: string, url?: string) : Observable<T>;
    // postFile<T>(query: string, fileName: string, file: any, url?: string) : Observable<T>;
}

export class DefaultApiService implements ApiService {

    private baseUrl: string = config.api.serverUrl;
    protected static authenticated$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);//localStorage.getItem(KEY_AUTH_TOKEN)?.length > 0);

    private readonly DEFAULT_ERROR = {'origin': 'server', 'message': 'Server error'};
    private readonly FORBIDDEN_ERROR = {'origin': 'server', 'message': 'Access Denied'};
    private readonly UNAUTHORIZED_ERROR = {'origin': 'server', 'message': 'Unauthorized'};

    private static getLanguage(): string {
        return "de_DE";
    }

    public isAuthenticated(): Observable<boolean> {
        return DefaultApiService.authenticated$;
    }
    public setAuthenticated(value: boolean, token: string) {
        if (value && token) {
            localStorage.setItem(KEY_AUTH_TOKEN, token);
        } else {
            localStorage.setItem(KEY_AUTH_TOKEN, "");
        }
        DefaultApiService.authenticated$.next(value);
    }

    private prepareError(err: any) : any {
        if (err && err.response) {
            const status = err.response.status;
            if (status === 403) {
                return this.FORBIDDEN_ERROR;
            } else if (status === 401) {
                DefaultApiService.authenticated$.next(false);
                return this.UNAUTHORIZED_ERROR;
            }
        }
        return err || this.DEFAULT_ERROR;
    }

    private static getOption(options: any, key: any, defaultValue: any): string {
        if (options) {
            if (options.hasOwnProperty(key)) {
                return options[key];
            }
        }
        return defaultValue;
    }

    protected getHeaders(options?: {}): any {
        let headers: any = {};
        const token = localStorage.getItem(KEY_AUTH_TOKEN);
        if (token) {
            headers[HEADER_AUTHORIZATION] = 'Bearer ' + token;
        }
        let language = DefaultApiService.getLanguage();
        if (language) {
            headers[HEADER_LANGUAGE] = language;
        }
        let value = DefaultApiService.getOption(options, HEADER_CONTENT_TYPE, 'application/json; charset=utf-8');
        if (value) {
            headers[HEADER_CONTENT_TYPE] = value;
        }
        headers[HEADER_ACCEPT] = 'application/json';
        return headers;
    }

    protected getUrl(query: string, url?: string) {
        return (url ? url : this.baseUrl) + query;
    }

    get<T>(query: string, url?: string): Observable<T> {
        return new Observable((observer) => {
            axios.get(this.getUrl(query, url), {headers: this.getHeaders()})
                .then((response) => {
                    observer.next(response.data);
                    observer.complete();
                })
                .catch((error) => observer.error(this.prepareError(error)));
        });
    }

    post<T>(query: string, payload: any, url?: string): Observable<T> {
        return new Observable((observer) => {
            axios.post<T>(this.getUrl(query, url), payload, {headers: this.getHeaders()})
                .then((response) => {
                    observer.next(response.data);
                    observer.complete();
                })
                .catch((error) => observer.error(this.prepareError(error)));
        });
    }

    put<T>(query: string, payload: any, url?: string): Observable<T> {
        return new Observable((observer) => {
            axios.put<T>(this.getUrl(query, url), payload, {headers: this.getHeaders()})
                .then((response) => {
                    observer.next(response.data);
                    observer.complete();
                })
                .catch((error) => observer.error(this.prepareError(error)));
        });
    }

    delete<T>(query: string, url?: string) : Observable<T> {
        return new Observable((observer) => {
            axios.delete(this.getUrl(query, url), {headers: this.getHeaders()})
                .then((response) => {
                    observer.next(response.data);
                    observer.complete();
                })
                .catch((error) => observer.error(this.prepareError(error)));
        });
    }

    postFile<T>(query: string, fileName: string, file: any, url?: string) : Observable<T> {
        let fd = new FormData();
        fd.append("fileName", fileName);
        fd.append("file", file, fileName);

        return new Observable((observer) => {
            axios.post<T>(this.getUrl(query, url), fd, {headers: this.getHeaders({[HEADER_CONTENT_TYPE]: undefined})})
                .then((response) => {
                    observer.next(response.data);
                    observer.complete();
                })
                .catch((error) => observer.error(this.prepareError(error)));
        });
    }

}
