// ComboDetailPage — Combo detail with flight+hotel+tour breakdown
function ComboDetailPage({ navigate, params }) {
  const combo = params?.combo || COMBOS[0];
  const [activeTab, setActiveTab] = useState('overview');
  const [selDep, setSelDep] = useState(null);
  const [adults, setAdults] = useState(2);
  const [children, setChildren] = useState(0);
  const [voucher, setVoucher] = useState('');
  const [vOk, setVOk] = useState(false);
  const [saved, setSaved] = useState(false);

  const dep = DEPARTURES.find(d => d.id === selDep);
  const basePrice = dep ? dep.priceAdult : Math.round(combo.price / 2);
  const childPrice = Math.round(basePrice * 0.7);
  const subtotal = adults * basePrice + children * childPrice;
  const disc = vOk ? Math.round(subtotal * 0.1) : 0;
  const finalTotal = subtotal - disc;

  // Combo component data
  const isBeach = ['Phú Quốc','Đà Nẵng','Nha Trang'].includes(combo.dest);
  const hasFlight = combo.id !== 3;
  const flightData = {
    airline: combo.id === 1 ? 'VietJet Air' : combo.id === 2 ? 'Vietnam Airlines' : null,
    go: { from: combo.from.split(' / ')[0], to: combo.dest, time: '07:30 → 09:15', flight: combo.id===1?'VJ148':'VN238' },
    back: { from: combo.dest, to: combo.from.split(' / ')[0], time: '20:30 → 22:15', flight: combo.id===1?'VJ149':'VN239' },
    class: 'Economy · 23kg hành lý ký gửi',
  };
  const transportData = {
    type: 'Xe limousine cao cấp 9 chỗ',
    route: combo.from + ' ↔ ' + combo.dest,
    time: '~8 tiếng / chiều',
  };
  const hotelData = {
    name: combo.id===1 ? 'JW Marriott Phu Quoc Emerald Bay' : combo.id===2 ? 'Fusion Suites Da Nang Beach' : 'Sapa House Hotel & Spa',
    stars: combo.id===1 ? 5 : 4,
    nights: parseInt(combo.duration[0]) - 1,
    breakfast: true,
    pool: isBeach,
  };
  const tourData = combo.id !== 2 ? {
    name: combo.id===1 ? 'Tour khám phá đảo Phú Quốc full day' : 'Trekking bản Cát Cát & Chinh phục Fansipan',
    duration: combo.id===1 ? '1 ngày' : '2 ngày',
    guide: true,
  } : null;

  const tabs = [
    { id:'overview', label:'Mô tả' },
    { id:'components', label:'Thành phần' },
    { id:'itinerary', label:'Lịch trình' },
    { id:'policy', label:'Chính sách' },
    { id:'reviews', label:'Đánh giá' },
  ];

  return (
    <div className="page-anim" style={{ paddingTop:68 }}>
      {/* Breadcrumb */}
      <div style={{ borderBottom:'1px solid var(--border-lt)', background:'var(--surface)' }}>
        <div className="container" style={{ padding:'12px 32px', display:'flex', gap:6, fontSize:13, color:'var(--text3)', alignItems:'center' }}>
          <button onClick={() => navigate('home')} style={{ background:'none',border:'none',cursor:'pointer',color:'var(--primary)',fontSize:13,padding:0,fontWeight:500 }}>Trang chủ</button>
          <Ico.chevR />
          <button onClick={() => navigate('search',{type:'combo'})} style={{ background:'none',border:'none',cursor:'pointer',color:'var(--primary)',fontSize:13,padding:0,fontWeight:500 }}>Combo</button>
          <Ico.chevR />
          <span style={{ color:'var(--text)', fontWeight:500 }}>{combo.dest} {combo.duration}</span>
        </div>
      </div>

      <div className="container" style={{ paddingTop:28, paddingBottom:48 }}>
        {/* Title */}
        <div style={{ marginBottom:16 }}>
          <div style={{ display:'flex', gap:8, marginBottom:10, flexWrap:'wrap' }}>
            <span className="badge badge-danger" style={{ fontSize:12, padding:'4px 12px' }}>Tiết kiệm {combo.disc}%</span>
            <span className="badge badge-primary">{combo.duration}</span>
            {hasFlight && <span className="badge badge-neutral" style={{ display:'flex', gap:3, alignItems:'center' }}><Ico.plane />Bay nội địa</span>}
            <span className="badge badge-neutral">Từ: {combo.from}</span>
          </div>
          <h1 style={{ fontSize:28, fontWeight:900, letterSpacing:'-.02em', lineHeight:1.2, marginBottom:10 }}>{combo.name}</h1>
          <div style={{ display:'flex', alignItems:'center', gap:16, flexWrap:'wrap' }}>
            <RatingRow score={4.8} count={187} size={14} />
            <div style={{ display:'flex', gap:8, marginLeft:'auto' }}>
              <button onClick={() => setSaved(!saved)} className="btn btn-ghost btn-sm" style={{ gap:5, color:saved?'var(--danger)':'var(--text2)' }}>
                <Ico.heart filled={saved} />{saved?'Đã lưu':'Lưu combo'}
              </button>
              <button className="btn btn-ghost btn-sm" style={{ gap:5 }}><Ico.share />Chia sẻ</button>
            </div>
          </div>
        </div>

        {/* Gallery */}
        <div style={{ display:'grid', gridTemplateColumns:'2fr 1fr', gridTemplateRows:'200px 200px', gap:10, borderRadius:'var(--r-xl)', overflow:'hidden', marginBottom:32 }}>
          <Img bg={combo.bg} stripe={combo.stripe} label={combo.label} style={{ gridRow:'1 / 3' }} />
          <Img bg="#0e7490" stripe="#0a5570" label={hasFlight ? flightData.airline + '\n' + flightData.go.from + ' → ' + flightData.go.to : 'Xe limousine'} style={{ height:'100%' }} />
          <Img bg="#2d6a4f" stripe="#1b4332" label={hotelData.name} style={{ height:'100%' }} />
        </div>

        {/* Savings highlight bar */}
        <div style={{ background:'linear-gradient(135deg, var(--accent) 0%, oklch(0.55 0.22 20) 100%)', borderRadius:'var(--r-lg)', padding:'16px 24px', marginBottom:28, display:'flex', alignItems:'center', justifyContent:'space-between', flexWrap:'wrap', gap:12 }}>
          <div style={{ color:'#fff' }}>
            <div style={{ fontSize:13, fontWeight:600, opacity:.8, marginBottom:2 }}>Tiết kiệm so với đặt riêng lẻ</div>
            <div style={{ fontSize:24, fontWeight:900 }}>{formatPriceFull(combo.oldPrice - combo.price)}</div>
          </div>
          <div style={{ display:'flex', gap:16 }}>
            <div style={{ textAlign:'center', color:'#fff' }}>
              <div style={{ fontSize:11, opacity:.7 }}>Giá gốc</div>
              <div style={{ fontWeight:700, textDecoration:'line-through', opacity:.7 }}>{formatPrice(combo.oldPrice)}</div>
            </div>
            <div style={{ textAlign:'center', color:'#fff' }}>
              <div style={{ fontSize:11, opacity:.7 }}>Giá combo</div>
              <div style={{ fontSize:20, fontWeight:900 }}>{formatPrice(combo.price)}</div>
            </div>
          </div>
        </div>

        {/* Layout */}
        <div style={{ display:'grid', gridTemplateColumns:'1fr 348px', gap:28, alignItems:'flex-start' }}>
          {/* Content */}
          <div>
            <div style={{ display:'flex', gap:0, borderBottom:'2px solid var(--border-lt)', marginBottom:28 }}>
              {tabs.map(t => (
                <button key={t.id} onClick={() => setActiveTab(t.id)}
                  style={{ padding:'10px 18px', border:'none', cursor:'pointer', fontWeight:700, fontSize:13, background:'transparent', transition:'all .15s', borderBottom: activeTab===t.id ? '3px solid var(--primary)' : '3px solid transparent', color: activeTab===t.id ? 'var(--primary)' : 'var(--text2)', marginBottom:'-2px' }}>
                  {t.label}
                </button>
              ))}
            </div>
            {activeTab === 'overview' && <ComboOverviewTab combo={combo} />}
            {activeTab === 'components' && <ComboComponentsTab hasFlight={hasFlight} flightData={flightData} transportData={transportData} hotelData={hotelData} tourData={tourData} />}
            {activeTab === 'itinerary' && <ComboItineraryTab combo={combo} hotelData={hotelData} hasFlight={hasFlight} />}
            {activeTab === 'policy' && <PolicyTab />}
            {activeTab === 'reviews' && <ReviewsTab tour={{ rating:4.8, reviews:187 }} />}
          </div>

          {/* Booking Sidebar */}
          <div style={{ position:'sticky', top:84 }}>
            <div className="card-flat" style={{ borderRadius:'var(--r-xl)', padding:'24px', border:'2px solid var(--border)' }}>
              <div style={{ marginBottom:14 }}>
                <div style={{ fontSize:12, color:'var(--text3)', textDecoration:'line-through' }}>{formatPriceFull(combo.oldPrice)}</div>
                <div style={{ fontSize:26, fontWeight:900, color:'var(--accent)' }}>{formatPriceFull(combo.price)}</div>
                <div style={{ fontSize:13, color:'var(--success)', fontWeight:600 }}>Tiết kiệm {combo.disc}% so với đặt riêng</div>
              </div>
              <hr style={{ marginBottom:16 }} />

              {/* Departure selection */}
              <div style={{ marginBottom:16 }}>
                <div style={{ fontWeight:700, fontSize:14, marginBottom:10 }}>Chọn đợt khởi hành</div>
                <div style={{ display:'flex', flexDirection:'column', gap:8 }}>
                  {DEPARTURES.filter(d => d.booked < d.slots).slice(0,4).map(d => (
                    <button key={d.id} onClick={() => setSelDep(d.id)}
                      style={{ padding:'10px 14px', borderRadius:'var(--r-md)', border: selDep===d.id ? '2px solid var(--primary)' : '1.5px solid var(--border)', background: selDep===d.id ? 'var(--primary-xlt)' : '#fff', cursor:'pointer', textAlign:'left', transition:'all .15s', display:'flex', justifyContent:'space-between', alignItems:'center' }}>
                      <div>
                        <span style={{ fontWeight:700, fontSize:14, color: selDep===d.id ? 'var(--primary)' : 'var(--text)' }}>{d.date}</span>
                        <span style={{ fontSize:12, color:'var(--text3)', marginLeft:6 }}>{d.dayOfWeek}</span>
                      </div>
                      <span className={'badge badge-' + (d.slots-d.booked<=5 ? 'warn' : 'success')}>Còn {d.slots-d.booked}</span>
                    </button>
                  ))}
                </div>
              </div>

              {/* Passengers */}
              <div style={{ marginBottom:16 }}>
                <div style={{ fontWeight:700, fontSize:14, marginBottom:10 }}>Số hành khách</div>
                <div style={{ display:'flex', flexDirection:'column', gap:10, padding:'14px', background:'var(--bg)', borderRadius:'var(--r-md)' }}>
                  <NumStepper label="Người lớn" sublabel="≥ 12 tuổi" value={adults} onChange={setAdults} min={1} max={10} />
                  <NumStepper label="Trẻ em" sublabel="5–11 tuổi · 70% giá" value={children} onChange={setChildren} min={0} max={5} />
                </div>
              </div>

              {/* Voucher */}
              <div style={{ marginBottom:16 }}>
                <div style={{ display:'flex', gap:8 }}>
                  <input className="input" placeholder="Mã giảm giá..." value={voucher} onChange={e => setVoucher(e.target.value.toUpperCase())} style={{ flex:1, fontSize:13 }} />
                  <button onClick={() => voucher && setVOk(true)} className="btn btn-outline btn-sm">Áp dụng</button>
                </div>
                {vOk && <div style={{ marginTop:6, fontSize:12, color:'var(--success)', fontWeight:600, display:'flex', gap:4, alignItems:'center' }}><Ico.check />Giảm 10%</div>}
              </div>

              {/* Price breakdown */}
              <div style={{ background:'var(--bg)', borderRadius:'var(--r-md)', padding:'14px', marginBottom:16 }}>
                <div style={{ display:'flex', justifyContent:'space-between', fontSize:13, color:'var(--text2)', marginBottom:6 }}>
                  <span>Người lớn × {adults}</span><span>{formatPriceFull(adults * basePrice)}</span>
                </div>
                {children > 0 && <div style={{ display:'flex', justifyContent:'space-between', fontSize:13, color:'var(--text2)', marginBottom:6 }}>
                  <span>Trẻ em × {children}</span><span>{formatPriceFull(children * childPrice)}</span>
                </div>}
                {disc > 0 && <div style={{ display:'flex', justifyContent:'space-between', fontSize:13, color:'var(--success)', marginBottom:6 }}>
                  <span>Voucher</span><span>−{formatPriceFull(disc)}</span>
                </div>}
                <hr style={{ margin:'8px 0' }} />
                <div style={{ display:'flex', justifyContent:'space-between', fontWeight:800, fontSize:16 }}>
                  <span>Tổng cộng</span><span style={{ color:'var(--accent)' }}>{formatPriceFull(finalTotal)}</span>
                </div>
              </div>

              <button onClick={() => navigate('booking', { combo, adults, children, dep, type:'combo' })}
                className="btn btn-accent btn-lg" style={{ width:'100%', borderRadius:'var(--r-md)', justifyContent:'center' }}>
                Đặt combo — {formatPrice(finalTotal)}
              </button>
              <button className="btn btn-ghost" style={{ width:'100%', marginTop:8, justifyContent:'center', fontSize:13 }}>Tư vấn miễn phí</button>
              <div style={{ marginTop:14, display:'flex', flexDirection:'column', gap:6 }}>
                {['Bao gồm ' + (hasFlight ? 'vé máy bay' : 'xe limousine'), 'Bao gồm khách sạn ' + '★'.repeat(hotelData.stars), tourData ? 'Bao gồm tour tham quan' : 'Ăn sáng tại khách sạn'].map(t => (
                  <div key={t} style={{ fontSize:12, color:'var(--success)', display:'flex', gap:5, alignItems:'center' }}><Ico.check />{t}</div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function ComboOverviewTab({ combo }) {
  return (
    <div>
      <div style={{ marginBottom:24 }}>
        <h3 style={{ fontWeight:800, fontSize:18, marginBottom:12 }}>Mô tả combo</h3>
        <p style={{ fontSize:14, lineHeight:1.8, color:'var(--text2)' }}>
          Combo {combo.dest} {combo.duration} là giải pháp du lịch hoàn hảo dành cho bạn — gói trọn gói tất cả những gì bạn cần trong một lần đặt. Thay vì mất hàng giờ so sánh và tự sắp xếp từng dịch vụ, bạn có ngay một hành trình được thiết kế chuyên nghiệp với giá tốt nhất.
        </p>
        <p style={{ fontSize:14, lineHeight:1.8, color:'var(--text2)', marginTop:10 }}>
          Đây là combo được khách hàng đặt nhiều nhất trong tháng 6/2026, với hơn 187 lượt đặt thành công và đánh giá 4.8/5 sao từ khách hàng thực tế.
        </p>
      </div>
      <div style={{ marginBottom:24 }}>
        <h3 style={{ fontWeight:800, fontSize:18, marginBottom:14 }}>Combo bao gồm</h3>
        <div style={{ display:'flex', flexDirection:'column', gap:8 }}>
          {combo.includes.map((inc,i) => (
            <div key={i} style={{ display:'flex', gap:12, alignItems:'center', padding:'12px 16px', background:'var(--bg)', borderRadius:'var(--r-md)' }}>
              <span style={{ color:'var(--success)', flexShrink:0 }}><Ico.check /></span>
              <span style={{ fontSize:14, fontWeight:500 }}>{inc}</span>
            </div>
          ))}
        </div>
      </div>
      <div style={{ padding:'16px 20px', background:'var(--warn-lt)', borderRadius:'var(--r-md)', fontSize:13, color:'var(--text2)', lineHeight:1.6 }}>
        <strong style={{ color:'var(--warn)' }}>Lưu ý quan trọng:</strong> Vé máy bay trong combo là thông tin tham khảo về hãng bay và chặng. Lịch bay cụ thể (giờ bay, số hiệu chuyến) sẽ được đại lý xác nhận sau khi booking. Đại lý chịu trách nhiệm đặt vé thực tế và thông báo cho khách.
      </div>
    </div>
  );
}

function ComboComponentsTab({ hasFlight, flightData, transportData, hotelData, tourData }) {
  const compCard = (icon, title, badge, children) => (
    <div style={{ border:'1.5px solid var(--border)', borderRadius:'var(--r-lg)', padding:'20px', background:'#fff' }}>
      <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:16 }}>
        <div style={{ fontWeight:800, fontSize:16, display:'flex', gap:8, alignItems:'center' }}>
          <span style={{ fontSize:22 }}>{icon}</span> {title}
        </div>
        <span className="badge badge-primary">{badge}</span>
      </div>
      {children}
    </div>
  );

  return (
    <div style={{ display:'flex', flexDirection:'column', gap:16 }}>
      {/* Flight or Transport */}
      {compCard(hasFlight ? '✈' : '🚌', hasFlight ? flightData.airline : 'Phương tiện di chuyển', 'Bao gồm',
        hasFlight ? (
          <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:12 }}>
            {[
              { lbl:'Chiều đi', route:`${flightData.go.from} → ${flightData.go.to}`, time:flightData.go.time, flight:flightData.go.flight },
              { lbl:'Chiều về', route:`${flightData.back.from} → ${flightData.back.to}`, time:flightData.back.time, flight:flightData.back.flight },
            ].map(f => (
              <div key={f.lbl} style={{ padding:'14px', background:'var(--primary-xlt)', borderRadius:'var(--r-md)' }}>
                <div style={{ fontSize:11, color:'var(--text3)', fontWeight:700, marginBottom:4, textTransform:'uppercase' }}>{f.lbl}</div>
                <div style={{ fontWeight:800, fontSize:16, color:'var(--primary)', marginBottom:3 }}>{f.route}</div>
                <div style={{ fontSize:13, color:'var(--text2)' }}>{f.time}</div>
                <div style={{ fontSize:11, color:'var(--text3)', marginTop:3 }}>{f.flight}</div>
              </div>
            ))}
            <div style={{ gridColumn:'1/-1', fontSize:12, color:'var(--text3)', display:'flex', gap:6, alignItems:'center' }}>
              <Ico.check />{flightData.class}
            </div>
          </div>
        ) : (
          <div style={{ padding:'14px', background:'var(--primary-xlt)', borderRadius:'var(--r-md)' }}>
            <div style={{ fontWeight:700, fontSize:15, marginBottom:4 }}>{transportData.type}</div>
            <div style={{ fontSize:13, color:'var(--text2)' }}>{transportData.route}</div>
            <div style={{ fontSize:12, color:'var(--text3)', marginTop:4 }}>Thời gian: {transportData.time}</div>
          </div>
        )
      )}

      {/* Hotel */}
      {compCard('🏨', 'Khách sạn & Lưu trú', 'Bao gồm',
        <div style={{ display:'flex', gap:16, alignItems:'center' }}>
          <Img bg={hasFlight ? '#0369a1' : '#2d6a4f'} stripe={hasFlight ? '#024d79' : '#1b4332'} label="" style={{ width:90, height:70, borderRadius:'var(--r-sm)', flexShrink:0 }} />
          <div style={{ flex:1 }}>
            <div style={{ fontWeight:800, fontSize:15, marginBottom:4 }}>{hotelData.name}</div>
            <div style={{ color:'var(--gold)', fontSize:16, marginBottom:4 }}>{'★'.repeat(hotelData.stars)}</div>
            <div style={{ display:'flex', gap:10, fontSize:13, color:'var(--text2)', flexWrap:'wrap' }}>
              <span>{hotelData.nights} đêm</span>
              {hotelData.breakfast && <span style={{ color:'var(--success)', display:'flex', gap:3, alignItems:'center' }}><Ico.check />Ăn sáng</span>}
              {hotelData.pool && <span style={{ color:'var(--success)', display:'flex', gap:3, alignItems:'center' }}><Ico.check />Hồ bơi</span>}
            </div>
          </div>
        </div>
      )}

      {/* Tour (if applicable) */}
      {tourData && compCard('🗺', 'Tour tham quan', 'Bao gồm',
        <div style={{ padding:'14px', background:'var(--bg)', borderRadius:'var(--r-md)' }}>
          <div style={{ fontWeight:700, fontSize:15, marginBottom:4 }}>{tourData.name}</div>
          <div style={{ display:'flex', gap:10, fontSize:13, color:'var(--text2)' }}>
            <span>{tourData.duration}</span>
            {tourData.guide && <span style={{ color:'var(--success)', display:'flex', gap:3, alignItems:'center' }}><Ico.check />Hướng dẫn viên</span>}
          </div>
        </div>
      )}
    </div>
  );
}

function ComboItineraryTab({ combo, hotelData, hasFlight }) {
  const days = combo.id === 1
    ? [
        { day:'Ngày 1', title:`${combo.from.split(' / ')[0]} → ${combo.dest}`, desc:'Bay đến Phú Quốc, nhận phòng resort, tự do khám phá bãi biển. Tối dạo phố Dương Đông, thưởng thức hải sản tươi ngon.', meals:['Tối'] },
        { day:'Ngày 2', title:'Khám phá đảo Phú Quốc', desc:'Tour đảo full day: Hòn Thơm, Hòn Gầm, lặn ngắm san hô, câu cá. Buffet trưa trên tàu. Chiều tự do tắm biển, massage.', meals:['Sáng','Trưa'] },
        { day:'Ngày 3', title:`${combo.dest} → ${combo.from.split(' / ')[0]}`, desc:'Ăn sáng, check-out, mua sắm đặc sản. Bay về. Kết thúc hành trình.', meals:['Sáng'] },
      ]
    : combo.id === 2
    ? [
        { day:'Ngày 1', title:`${combo.from.split(' / ')[0]} → Đà Nẵng`, desc:'Bay đến Đà Nẵng, nhận phòng khách sạn bãi biển Mỹ Khê. Buổi chiều tắm biển, buổi tối tham quan Cầu Rồng.', meals:['Tối'] },
        { day:'Ngày 2', title:'Bà Nà Hills — Cầu Vàng', desc:'Cáp treo lên Bà Nà Hills, khám phá Làng Pháp và Cầu Vàng nổi tiếng. Tham quan vườn hoa Le Jardin D\'Amour.', meals:['Sáng','Trưa'] },
        { day:'Ngày 3', title:'Hội An cổ kính', desc:'Xe đến phố cổ Hội An, tham quan chùa Cầu, hội quán Phúc Kiến, thuyền thúng làng rau Trà Quế. Tối ngắm đèn lồng trên sông Hoài.', meals:['Sáng'] },
        { day:'Ngày 4', title:'Đà Nẵng → Về nhà', desc:'Ăn sáng, tự do mua sắm đặc sản (bánh tráng cuốn thịt heo, nước mắm Nam Ô). Bay về.', meals:['Sáng'] },
      ]
    : [
        { day:'Ngày 1', title:`Hà Nội → Sapa`, desc:'Xe limousine đón tại Hà Nội 07:00, đến Sapa 14:00. Nhận phòng khách sạn view núi. Dạo chợ Sapa, ngắm hoàng hôn.', meals:['Tối'] },
        { day:'Ngày 2', title:'Trekking Cát Cát — Fansipan', desc:'Trekking bản Cát Cát, thăm làng H\'Mông. Chiều bắt cáp treo lên đỉnh Fansipan 3.147m. Ngắm mây và núi rừng hùng vĩ.', meals:['Sáng','Trưa','Tối'] },
        { day:'Ngày 3', title:'Sapa → Hà Nội', desc:'Ăn sáng, trekking thung lũng Mường Hoa ngắm ruộng bậc thang. 12:00 xe khởi hành về Hà Nội. 18:00 đến nơi.', meals:['Sáng'] },
      ];

  return (
    <div style={{ display:'flex', flexDirection:'column', gap:20 }}>
      {days.map((d,i) => (
        <div key={i} style={{ display:'flex', gap:20 }}>
          <div style={{ flexShrink:0, display:'flex', flexDirection:'column', alignItems:'center' }}>
            <div style={{ width:44, height:44, borderRadius:'50%', background:'var(--primary)', color:'#fff', display:'flex', alignItems:'center', justifyContent:'center', fontWeight:800, fontSize:12, textAlign:'center', lineHeight:1.2 }}>N{i+1}</div>
            {i < days.length-1 && <div style={{ width:2, flex:1, background:'var(--border-lt)', margin:'6px 0' }} />}
          </div>
          <div style={{ flex:1, paddingBottom: i < days.length-1 ? 16 : 0 }}>
            <div style={{ fontWeight:800, fontSize:16, marginBottom:6 }}>{d.day}: {d.title}</div>
            <div style={{ display:'flex', gap:6, marginBottom:10, flexWrap:'wrap' }}>
              {d.meals.map(m => <span key={m} style={{ fontSize:12, padding:'2px 9px', borderRadius:20, background:'var(--gold-lt)', color:'var(--gold)', fontWeight:600 }}>🍽 Bữa {m}</span>)}
            </div>
            <p style={{ fontSize:14, lineHeight:1.7, color:'var(--text2)' }}>{d.desc}</p>
          </div>
        </div>
      ))}
    </div>
  );
}

Object.assign(window, { ComboDetailPage });
