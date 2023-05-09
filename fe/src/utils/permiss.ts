export const setPermissionCodes = (codes: string[] | null) => {
    if (!codes) {
        localStorage.removeItem('permissionCodes');
        return;
    }
    localStorage.setItem('permissionCodes', JSON.stringify(codes));
}

export const getPermissionCodes = (): string[] | null => { 
    const codes = localStorage.getItem('permissionCodes');
    if (codes) {
        return JSON.parse(codes);
    }
    return null;
}