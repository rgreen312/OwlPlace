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