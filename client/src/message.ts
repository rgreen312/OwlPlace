export enum MsgType {
    ERROR = -1,
    OPEN = 0,
    DRAWPIXEL = 1,
    LOGINUSER = 2,
    CHANGECLIENTPIXEL = 3,
    IMAGE = 4,
    TESTING = 5,
    DRAWRESPONSE = 6,
    CLOSE = 9,
}

export interface Msg {
    type: number;
}

export class ErrorMsg implements Msg {
    type: number = MsgType.ERROR;
}

export class ImageMsg implements Msg {
    type: number = MsgType.IMAGE;
    formatString: string;
    
    constructor(formatString: string) {
        this.formatString = formatString;
    }
}

export class ChangeClientPixelMsg implements Msg {
    type: number = MsgType.CHANGECLIENTPIXEL;
    r: number;
    g: number;
    b: number;
    x: number;
    y: number;

    constructor(r: number, g: number, b: number, x: number, y: number) {
        this.r = r;
        this.g = g;
        this.b = b;
        this.x = x;
        this.y = y;

    }
    
}

function parseMsg(json : string) : Msg {
    let data = JSON.parse(json);
    switch (data.type) {
        case MsgType.IMAGE: {
            return new ImageMsg(data.formatString);
            break;
        }
        case MsgType.CHANGECLIENTPIXEL: {
            return new ChangeClientPixelMsg(data.r, data.g, data.b, data.x, data.y);
            break;
        }
        default: {
            return new ErrorMsg();
        }
    }
}