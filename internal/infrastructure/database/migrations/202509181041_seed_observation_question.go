package migrations

import (
	"backend-golang/internal/infrastructure/database/models"
	"time"

	"gorm.io/gorm"
)

func SeedObservationQuestionsTableUp(tx *gorm.DB) error {
	questions := []models.ObservationQuestion{
		{QuestionCode: "BPE-01", AgeCategory: "Balita", QuestionNumber: 1, QuestionText: "Hipoaktif atau bergerak tidak bertujuan", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BPE-02", AgeCategory: "Balita", QuestionNumber: 2, QuestionText: "Hipoaktif atau lamban gerak", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BPE-03", AgeCategory: "Balita", QuestionNumber: 3, QuestionText: "Tidak mampu mengikuti aturan", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BPE-04", AgeCategory: "Balita", QuestionNumber: 4, QuestionText: "Menyakiti diri sendiri", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BPE-05", AgeCategory: "Balita", QuestionNumber: 5, QuestionText: "Menyerang orang lain ketika marah", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BPE-06", AgeCategory: "Balita", QuestionNumber: 6, QuestionText: "Perilaku repetitif atau berulang-ulang", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BPE-07", AgeCategory: "Balita", QuestionNumber: 7, QuestionText: "Tidak dapat duduk tenang", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BPE-08", AgeCategory: "Balita", QuestionNumber: 8, QuestionText: "Anak Jalan jinjit", Score: 2, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "BFM-01", AgeCategory: "Balita", QuestionNumber: 9, QuestionText: "Kelainan pada anggota tubuh atau pemakaian alat bantu", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BFM-02", AgeCategory: "Balita", QuestionNumber: 10, QuestionText: "Tdak mampu melompat", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BFM-03", AgeCategory: "Balita", QuestionNumber: 11, QuestionText: "TIdak mampu mengikuti contoh gerakan seperti senam", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BFM-04", AgeCategory: "Balita", QuestionNumber: 12, QuestionText: "Tidak mampu membuat bentuk sederhana dari playdough", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BFM-05", AgeCategory: "Balita", QuestionNumber: 13, QuestionText: "Tidak mampu merobek kertas", Score: 2, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "BBB-01", AgeCategory: "Balita", QuestionNumber: 14, QuestionText: "Saat ditanya mengulang pertanyaan atau perkataan", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BBB-02", AgeCategory: "Balita", QuestionNumber: 15, QuestionText: "Tdak mampu memahami perintah/instruksi", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BBB-03", AgeCategory: "Balita", QuestionNumber: 16, QuestionText: "TIdak mampu berkomunikasi 2 arah.tanya jawab", Score: 3, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "BKA-01", AgeCategory: "Balita", QuestionNumber: 17, QuestionText: "Tidak mampu menyelesaikan aktifitas", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BKA-02", AgeCategory: "Balita", QuestionNumber: 18, QuestionText: "Tidak mampu mempertahankan atensi dan konsentrasi ketika diberi tugas", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BKA-03", AgeCategory: "Balita", QuestionNumber: 19, QuestionText: "Tidak mampu menyebutkan identitas diri dan anggota keluarga", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BKA-04", AgeCategory: "Balita", QuestionNumber: 20, QuestionText: "Tidak mampu menamai benda sekitar", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BKA-05", AgeCategory: "Balita", QuestionNumber: 21, QuestionText: "Tidak mampu menyebutkan angka 1-5", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BKA-06", AgeCategory: "Balita", QuestionNumber: 22, QuestionText: "Tidak mampu mengidentifikasi bentuk (minimal 1 bentuk konsisten)", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BKA-07", AgeCategory: "Balita", QuestionNumber: 23, QuestionText: "Tidak mampu mengidentifikasi warna primer", Score: 2, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "BS-01", AgeCategory: "Balita", QuestionNumber: 24, QuestionText: "Tidak ada kontak mata/kontak mata minim saat diajak berbicara", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BS-02", AgeCategory: "Balita", QuestionNumber: 25, QuestionText: "Suka menyendiri", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "BS-03", AgeCategory: "Balita", QuestionNumber: 26, QuestionText: "Kesulitan berdaptaasi dengan lingkungan baru", Score: 2, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "APE-01", AgeCategory: "Anak-anak", QuestionNumber: 1, QuestionText: "Hiperaktif atau bergerak tidak bertujuan", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "APE-02", AgeCategory: "Anak-anak", QuestionNumber: 2, QuestionText: "Hipoaktif atau lamban gerak", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "APE-03", AgeCategory: "Anak-anak", QuestionNumber: 3, QuestionText: "Tidak mampu mengikuti aturan", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "APE-04", AgeCategory: "Anak-anak", QuestionNumber: 4, QuestionText: "Menyakiti diri sendiri atau menyerang orang lain ketika marah", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "APE-05", AgeCategory: "Anak-anak", QuestionNumber: 5, QuestionText: "Perilaku Repetitif atau berulang-ulang", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "APE-06", AgeCategory: "Anak-anak", QuestionNumber: 6, QuestionText: "Tidak dapat duduk tenang", Score: 3, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "AFM-01", AgeCategory: "Anak-anak", QuestionNumber: 7, QuestionText: "Kelainan pada anggota tubuh atau pemakaian alat bantu", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AFM-02", AgeCategory: "Anak-anak", QuestionNumber: 8, QuestionText: "Tidak mampu melompat", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AFM-03", AgeCategory: "Anak-anak", QuestionNumber: 9, QuestionText: "Tidak mampu mengikuti contoh gerakan seperti senam", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AFM-04", AgeCategory: "Anak-anak", QuestionNumber: 10, QuestionText: "Tidak mampu menggunting", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AFM-05", AgeCategory: "Anak-anak", QuestionNumber: 11, QuestionText: "Tidak mampu melipat kertas", Score: 2, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "ABB-01", AgeCategory: "Anak-anak", QuestionNumber: 12, QuestionText: "Saat ditanya Mengulang pertanyaan atau perkataan", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "ABB-02", AgeCategory: "Anak-anak", QuestionNumber: 13, QuestionText: "Tidak mampu memahami perintah/instruksi", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "ABB-03", AgeCategory: "Anak-anak", QuestionNumber: 14, QuestionText: "Tidak mampu berkomunikasi 2 arah/tanya jawab", Score: 3, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "AKA-01", AgeCategory: "Anak-anak", QuestionNumber: 15, QuestionText: "Tidak mampu menyelesaikan tugas", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AKA-02", AgeCategory: "Anak-anak", QuestionNumber: 16, QuestionText: "Tidak mampu mempertahankan atensi dan konsentrasi ketika diberi tugas", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AKA-03", AgeCategory: "Anak-anak", QuestionNumber: 17, QuestionText: "Tidak mampu menyebutkan identitas diri dan anggota keluarga", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AKA-04", AgeCategory: "Anak-anak", QuestionNumber: 18, QuestionText: "Tidak mampu menamai benda sekitar", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AKA-05", AgeCategory: "Anak-anak", QuestionNumber: 19, QuestionText: "Tidak mampu mengurutkan angka 1-10", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AKA-06", AgeCategory: "Anak-anak", QuestionNumber: 20, QuestionText: "Tidak mampu mengurutkan abjad A-Z", Score: 1, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "AS-01", AgeCategory: "Anak-anak", QuestionNumber: 21, QuestionText: "Tidak ada kontak mata/kontak mata minim saat diajak berbicara", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AS-02", AgeCategory: "Anak-anak", QuestionNumber: 22, QuestionText: "Suka menyendiri", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AS-03", AgeCategory: "Anak-anak", QuestionNumber: 23, QuestionText: "Tidak mau berbagi dengan teman/egois", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "AS-04", AgeCategory: "Anak-anak", QuestionNumber: 24, QuestionText: "Kesulitan beradaptasi dengan lingkungan baru", Score: 2, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "PPE-01", AgeCategory: "Remaja", QuestionNumber: 1, QuestionText: "Hiperaktif atau bergerak tidak bertujuan", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "PPE-02", AgeCategory: "Remaja", QuestionNumber: 2, QuestionText: "Hipoaktif atau lamban gerak", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "PPE-03", AgeCategory: "Remaja", QuestionNumber: 3, QuestionText: "Tidak mampu mengikuti aturan", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "PPE-04", AgeCategory: "Remaja", QuestionNumber: 4, QuestionText: "Menyakiti diri sendiri atau menyerang orang lain", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "PPE-05", AgeCategory: "Remaja", QuestionNumber: 5, QuestionText: "Perilaku Repetitif atau berulang-ulang", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "PPE-06", AgeCategory: "Remaja", QuestionNumber: 6, QuestionText: "Tidak dapat duduk tenang", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "PPE-07", AgeCategory: "Remaja", QuestionNumber: 7, QuestionText: "Ketertarikan berlebih terhadap lawan jenis", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "PPE-08", AgeCategory: "Remaja", QuestionNumber: 8, QuestionText: "Emosi yang meledak-ledak", Score: 3, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "RFM-01", AgeCategory: "Remaja", QuestionNumber: 9, QuestionText: "Kelainan pada anggota tubuh atau pemakaian alat bantu", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RFM-02", AgeCategory: "Remaja", QuestionNumber: 10, QuestionText: "Tidak mampu menganyam", Score: 3, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "RBB-01", AgeCategory: "Remaja", QuestionNumber: 11, QuestionText: "Saat ditanya mengulang pertanyaan atau perkataan", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RBB-02", AgeCategory: "Remaja", QuestionNumber: 12, QuestionText: "Tidak mampu memahami perintah/instruksi tiga tahap", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RBB-03", AgeCategory: "Remaja", QuestionNumber: 13, QuestionText: "Tidak mampu berkomunikasi 2 arah/tanya jawab", Score: 3, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "RKA-01", AgeCategory: "Remaja", QuestionNumber: 14, QuestionText: "Tidak mampu menyelesaikan tugas", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RKA-02", AgeCategory: "Remaja", QuestionNumber: 15, QuestionText: "Tidak mampu mempertahankan atensi dan konsentrasi ketika diberi tugas", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RKA-03", AgeCategory: "Remaja", QuestionNumber: 16, QuestionText: "Tidak mampu menceritakan diri sendiri", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RKA-04", AgeCategory: "Remaja", QuestionNumber: 17, QuestionText: "Tidak mampu operasi hitung sederhana", Score: 2, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RKA-05", AgeCategory: "Remaja", QuestionNumber: 18, QuestionText: "Tidak mampu membaca paragraf sederhana", Score: 2, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "RS-01", AgeCategory: "Remaja", QuestionNumber: 19, QuestionText: "Tidak ada kontak mata/kontak mata minim saat diajak berbicara", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RS-02", AgeCategory: "Remaja", QuestionNumber: 20, QuestionText: "Suka menyendiri", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RS-03", AgeCategory: "Remaja", QuestionNumber: 21, QuestionText: "Kesulitan beradaptasi dengan lingkungan baru", Score: 2, IsActive: true, CreatedAt: time.Now()},

		{QuestionCode: "RK-01", AgeCategory: "Remaja", QuestionNumber: 22, QuestionText: "Tidak bisa mengancing baju sendiri", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RK-02", AgeCategory: "Remaja", QuestionNumber: 23, QuestionText: "Tidak bisa toilet training", Score: 3, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RK-03", AgeCategory: "Remaja", QuestionNumber: 24, QuestionText: "Tidak berpenampilan rapi dan sopan", Score: 1, IsActive: true, CreatedAt: time.Now()},
		{QuestionCode: "RK-04", AgeCategory: "Remaja", QuestionNumber: 25, QuestionText: "Tidak mengenal mata uang", Score: 2, IsActive: true, CreatedAt: time.Now()},
	}

	for _, question := range questions {
		var count int64
		if err := tx.Model(&models.ObservationQuestion{}).
			Where("question_code = ?", question.QuestionCode).
			Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {

			if err := tx.Create(&question).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func SeedObservationQuestionsTableDown(tx *gorm.DB) error {
	return tx.Where("question_code LIKE ?", "BPE-% OR question_code LIKE ? OR ...").
		Delete(&models.ObservationQuestion{}).Error
}
