export const ERROR = -1;
export const OPEN = 0;
export const DRAWPIXEL = 1;
export const LOGINUSER = 2;
export const UPDATEPIXEL = 3;
export const IMAGE = 4;
export const TESTING = 5;
export const DRAWRESPONSE = 6;
export const CLOSE = 9;

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