export const ERROR = -1;
export const IMAGE = 4;

export interface Msg {
    type: number;
}

export class ErrorMsg implements Msg {
    type: number = ERROR;
}

export class ImageMsg implements Msg {
    type: number = IMAGE;
    formatString: string;
    
    constructor(formatString: string) {
        this.formatString = formatString;
    }
}

function parseMsg(json : string) : Msg {
    let data = JSON.parse(json);
    switch (data.type) {
        case IMAGE: {
            return new ImageMsg(data.formatString);
            break;
        }
        default: {
            return new ErrorMsg();
        }
    }
}