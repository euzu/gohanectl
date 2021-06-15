
const formatDate = (value: string, duration?: boolean) => {
    if (value) {
        try {
            const now = new Date();
            const date = new Date(value);
            // let totalSeconds = (now.getTime() - date.getTime()) / 1000;
            // if (totalSeconds < 86400) {
            //     const hours = Math.floor(totalSeconds / 3600);
            //     totalSeconds %= 3600;
            //     const minutes = Math.floor(totalSeconds / 60);
            //     const seconds = Math.floor(totalSeconds % 60);
            //     let str = '';
            //     if (hours > 0) {
            //         str += hours + 'h ';
            //     }
            //     if (minutes > 0) {
            //         str += minutes + 'm ';
            //     }
            //     if (seconds > 0) {
            //         str += seconds + 's';
            //     }
            //     return str;
            // }
            const formatted = date.toLocaleString();
            if (now.getDate() === date.getDate() && now.getMonth() === date.getMonth() && now.getFullYear() === date.getFullYear()) {
                const idx = formatted.indexOf(',')
                return formatted.substring(idx+1);
            }
            return formatted.replace(',', '');
        } catch (e) {
        }
    }
    return '';
}

export default function useDateTimeFormat() {
    return formatDate;
}
