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
    VERIFICATIONFAIL = 10,
    CREATEUSER = 11
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
    
export class VerificationFailMsg implements Msg {
    type: number = MsgType.VERIFICATIONFAIL;
    status: number;
    
    constructor(status: number) {
        this.status = status;
    }
}

export class CreateUserMsg implements Msg {
    type: number = MsgType.CREATEUSER;
    status: number;
    cooldown: number;
    
    constructor(status: number, cooldown: number) {
        this.status = status;
        this.cooldown = cooldown;
    }
}
