INSERT INTO public.merchants (
    id,
    name,
    merchant_category,
    image_url,
    lat,
    long,
    created_at
) VALUES
('0199a00a-2888-7cbd-b4a6-6f8ad0950599', 'P9 Toko Buku Ilmu Jaya', 'SmallRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.553264, 95.316438, 1759582391413874000),
('019995ba-65b1-7749-9ddd-9bd6f7cceada', 'P2 Martabak Bang', 'SmallRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.553147, 95.318642, 1759582391413978000),
('01999619-405c-73e9-97da-a4b100b23465', 'P6 Kedai Mie Razali', 'SmallRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.553874, 95.316278, 1759582391413989000),
('019998b5-a219-7895-9aaa-3f0dee6448d1', 'P8 Minimarket Seulawah', 'SmallRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.554723, 95.317621, 1759582391413992000),
('01999616-1bed-7a80-87d7-75df3f5a4af2', 'P4 Warung Nasi Gurih', 'SmallRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.554612, 95.316842, 1759582391413994000),
('0199a00a-e101-7777-9774-6fafa4072185', 'P10 Restoran Seafood', 'MediumRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.552941, 95.318873, 1759582391413996000),
('019994ac-4344-7e6b-80f2-acb522774962', 'P1 Kopi Ulee Kareng', 'LargeRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.553982, 95.318274, 1759582391413998000),
('019995bd-0e72-7792-bb34-3f5b41718a18', 'P3 Roti Cane Simpang', 'MerchandiseRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.552743, 95.317089, 1759582391414001000),
('01999616-995c-761e-8f3f-2cd3083647cf', 'P5 Ayam Penyet Lampriet', 'ConvenienceStore', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.554976, 95.318425, 1759582391414003000),
('019998b4-ab91-74a0-8323-cdfe9b0b9533', 'P7 Apotek Sehat', 'BoothKiosk', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.552534, 95.317964, 1759582391414005000),
('0199a3e2-a0df-74e0-b90b-07ce6ef3bc16', 'Restoran Laut Lhoknga', 'SmallRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.475321, 95.247893, 1759582391414007000),
('0199a3e4-4baf-7ad1-9312-b33669f6e37b', 'Restoran Jantho Asri', 'SmallRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.434278, 95.541287, 1759582391414009000),
('0199a3e4-e695-725c-87eb-ad9cd360947b', 'Restoran Darussalam Indah', 'SmallRestaurant', 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 5.600842, 95.390756, 1759582391414011000);


INSERT INTO public.items (
    id,
    merchant_id,
    name,
    product_category,
    price,
    image_url,
    created_at
) VALUES
('019995b8-04db-77c3-bc94-389c22dbe06d', '019994ac-4344-7e6b-80f2-acb522774962', 'allopama', 'Beverage', 2000.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759153489116871000),
('0199997b-28d8-73d1-b29f-4419ee50e574', '019995ba-65b1-7749-9ddd-9bd6f7cceada', 'balwan', 'Condiments', 20.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759216609497047000),
('0199997b-6b31-7dfc-83ea-96fe1c1716bb', '019995bd-0e72-7792-bb34-3f5b41718a18', 'resol', 'Condiments', 30.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759216626482176000),
('0199997b-93de-7c45-bed6-acf51d8ca5c1', '01999616-1bed-7a80-87d7-75df3f5a4af2', 'mim', 'Beverage', 30.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759216636895074000),
('0199997b-cf32-76cd-9764-5214c5512f8b', '01999616-995c-761e-8f3f-2cd3083647cf', 'timphan', 'Food', 10.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759216652082657000),
('0199997b-f14e-7aa8-8989-4cb8a714c90b', '01999619-405c-73e9-97da-a4b100b23465', 'bhoi', 'Beverage', 10.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759216660814917000),
('0199997c-181d-7214-b125-2d52f81b8d1e', '019998b4-ab91-74a0-8323-cdfe9b0b9533', 'tape', 'Food', 110.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759216670749389000),
('0199997c-4eea-74b6-a82a-a6e745d74ce8', '019998b5-a219-7895-9aaa-3f0dee6448d1', 'karah', 'Additions', 10.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759216684778470000),
('0199a3e5-b9ef-7564-9694-ee559d4f2f15', '0199a3e4-e695-725c-87eb-ad9cd360947b', 'eret', 'Snack', 1.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759391365615875000),
('0199ae98-0c8d-79f1-9105-53f63627d2d3', '019994ac-4344-7e6b-80f2-acb522774962', 'Bika', 'Additions', 2000.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759570824334632000),
('0199af28-67c0-759f-aa8f-4cfc561b5ea7', '019995ba-65b1-7749-9ddd-9bd6f7cceada', 'Kjeu', 'Snack', 3000.00, 'https://bee.telkomuniversity.ac.id/wp-content/uploads/2024/10/solar-panels-8593759_1280.webp', 1759580284865319000);