import * as XLSX from 'xlsx'

export const buildExcel = (list,fileName) => {
    const data = XLSX.utils.json_to_sheet(list)
    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, data, 'Sheet1')
    
    const excelBuffer = XLSX.write(wb, { bookType: 'xlsx', type: 'array' });
    const blob = new Blob([excelBuffer], { type: 'application/octet-stream' });

    const name = fileName+'.xlsx';
    const link = document.createElement('a');
    link.href = window.URL.createObjectURL(blob);
    link.download = name;
    link.click();

}