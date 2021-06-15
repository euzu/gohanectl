import {Observable, Subject} from "rxjs";

let subject$: Subject<boolean> = new Subject<boolean>();
const INACTIVITY_THRESHOLD = 60000;
let lastEvent: number = 0;
const inactivity = () => lastEvent = performance.now();
window.addEventListener('click', inactivity);
window.addEventListener('scroll', inactivity);
window.addEventListener('keyup', inactivity);
window.addEventListener('blur', inactivity);
window.addEventListener('focus',function(){
    let diff = performance.now() - lastEvent;
    if (diff > INACTIVITY_THRESHOLD) {
        subject$.next(true);
    }
    lastEvent = performance.now();
});

export default function useWakeup() : Observable<boolean> {
    return subject$
}