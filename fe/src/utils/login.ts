import { setPermissionCodes } from './permiss';

export const setLoginInfo = (loginInfo: ILoginInfo | null) => {
    if (!loginInfo) {
        localStorage.removeItem('loginInfo');
        return;
    }
    localStorage.setItem('loginInfo', JSON.stringify(loginInfo));
}

export const getLoginInfo = (): ILoginInfo | null => {
    const loginInfo = localStorage.getItem('loginInfo');
    if (loginInfo) {
        return JSON.parse(loginInfo);
    }
    return null;
}


export const logout = () => {
    setLoginInfo(null);
    setPermissionCodes(null);
    window.location.href = "/";
}