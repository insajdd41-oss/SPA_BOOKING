
const API_URL = 'http://localhost:8080/api/bookings';

export async function fetchBookingsFromServer() {
  try {
    const response = await fetch(API_URL);
    
    if (!response.ok) {
      throw new Error(`Ошибка HTTP: ${response.status}`);
    }
    
    const data = await response.json();
    
    return data.map((item: any) => ({
      id: String(item.appointment_id),
      clientName: item.client_name,
      clientPhone: '+7 (не указан)', 
      serviceTitle: item.service_name,
      description: 'Описание услуги', 
      category: 'Массаж', 
      masterName: item.specialist_name,
      status: item.status === 'запланирована' ? 'new' : 'done', 
      price: item.price,
      durationMin: 60, 
      date: item.date.split('T')[0], 
      time: item.time,
      isPaid: false,
      createdAt: new Date().toISOString(),
    }));
    
  } catch (error) {
    console.error('Ошибка при загрузке записей с сервера:', error);
    return [];
  }
}