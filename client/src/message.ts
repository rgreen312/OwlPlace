export enum MsgType {
    ERROR = -1,
    OPEN = 0,
    DRAWPIXEL = 1,
    LOGINUSER = 2,
    UPDATEPIXEL = 3,
    IMAGE = 4,
    TESTING = 5,
    DRAWRESPONSE = 6,
    CLOSE = 9,
    VERIFICATIONFAIL = 10,
    CREATEUSER = 11
}

/*
    Error        MsgType = -1
	Open         MsgType = 0
	DrawPixel    MsgType = 1
    LoginUser    MsgType = 2
    UpdatePixel  MsgType = 3
	Image        MsgType = 4
	Testing      MsgType = 5
	DrawResponse MsgType = 6
    Close        MsgType = 9
    VerificationFail MsgType = 10
    CreateUser  MsgType = 11
    */

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
    
    constructor(status: number) {
        this.status = status;
    }
}

function parseMsg(json : string) : Msg {
    let data = JSON.parse(json);
    switch (data.type) {
        case MsgType.IMAGE: {
            return new ImageMsg(data.formatString);
            break;
        }
        default: {
            return new ErrorMsg();
        }
    }
}