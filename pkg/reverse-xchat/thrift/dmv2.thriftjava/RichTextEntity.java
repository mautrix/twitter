package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import androidx.media3.exoplayer.C8366s1;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000<\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\t\n\u0002\u0010\u000e\n\u0002\b\u0004\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 !2\u00020\u0001:\u0002\"!B%\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0005¢\u0006\u0004\b\u0007\u0010\bJ\u0017\u0010\f\u001a\u00020\u000b2\u0006\u0010\n\u001a\u00020\tH\u0016¢\u0006\u0004\b\f\u0010\rJ\u0012\u0010\u000e\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000e\u0010\u000fJ\u0012\u0010\u0010\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0010\u0010\u000fJ\u0012\u0010\u0011\u001a\u0004\u0018\u00010\u0005HÆ\u0003¢\u0006\u0004\b\u0011\u0010\u0012J4\u0010\u0013\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u0005HÆ\u0001¢\u0006\u0004\b\u0013\u0010\u0014J\u0010\u0010\u0016\u001a\u00020\u0015HÖ\u0001¢\u0006\u0004\b\u0016\u0010\u0017J\u0010\u0010\u0018\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0018\u0010\u0019J\u001a\u0010\u001d\u001a\u00020\u001c2\b\u0010\u001b\u001a\u0004\u0018\u00010\u001aHÖ\u0003¢\u0006\u0004\b\u001d\u0010\u001eR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001fR\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001fR\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00058\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u0010 ¨\u0006#"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextEntity;", "Lcom/bendb/thrifty/a;", "", "start_index", "end_index", "Lcom/x/dmv2/thriftjava/RichTextContent;", "content", "<init>", "(Ljava/lang/Integer;Ljava/lang/Integer;Lcom/x/dmv2/thriftjava/RichTextContent;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Integer;", "component2", "component3", "()Lcom/x/dmv2/thriftjava/RichTextContent;", "copy", "(Ljava/lang/Integer;Ljava/lang/Integer;Lcom/x/dmv2/thriftjava/RichTextContent;)Lcom/x/dmv2/thriftjava/RichTextEntity;", "", "toString", "()Ljava/lang/String;", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Integer;", "Lcom/x/dmv2/thriftjava/RichTextContent;", "Companion", "RichTextEntityAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class RichTextEntity implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final RichTextContent content;

    @JvmField
    @InterfaceC88465b
    public final Integer end_index;

    @JvmField
    @InterfaceC88465b
    public final Integer start_index;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new RichTextEntityAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/RichTextEntity$RichTextEntityAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/RichTextEntity;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/RichTextEntity;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/RichTextEntity;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class RichTextEntityAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public RichTextEntity m83780read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Integer numValueOf = null;
            Integer numValueOf2 = null;
            RichTextContent richTextContent = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new RichTextEntity(numValueOf, numValueOf2, richTextContent);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            C11272a.m14141a(protocol, b);
                        } else if (b == 12) {
                            richTextContent = (RichTextContent) RichTextContent.ADAPTER.read(protocol);
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 8) {
                        numValueOf2 = Integer.valueOf(protocol.mo14132c4());
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 8) {
                    numValueOf = Integer.valueOf(protocol.mo14132c4());
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a RichTextEntity struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("RichTextEntity");
            if (struct.start_index != null) {
                protocol.mo14136v3("start_index", 1, (byte) 8);
                protocol.mo14122C2(struct.start_index.intValue());
            }
            if (struct.end_index != null) {
                protocol.mo14136v3("end_index", 2, (byte) 8);
                protocol.mo14122C2(struct.end_index.intValue());
            }
            if (struct.content != null) {
                protocol.mo14136v3("content", 3, (byte) 12);
                RichTextContent.ADAPTER.write(protocol, struct.content);
            }
            protocol.mo14134i0();
        }
    }

    public RichTextEntity(@InterfaceC88465b Integer num, @InterfaceC88465b Integer num2, @InterfaceC88465b RichTextContent richTextContent) {
        this.start_index = num;
        this.end_index = num2;
        this.content = richTextContent;
    }

    public static /* synthetic */ RichTextEntity copy$default(RichTextEntity richTextEntity, Integer num, Integer num2, RichTextContent richTextContent, int i, Object obj) {
        if ((i & 1) != 0) {
            num = richTextEntity.start_index;
        }
        if ((i & 2) != 0) {
            num2 = richTextEntity.end_index;
        }
        if ((i & 4) != 0) {
            richTextContent = richTextEntity.content;
        }
        return richTextEntity.copy(num, num2, richTextContent);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Integer getStart_index() {
        return this.start_index;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final Integer getEnd_index() {
        return this.end_index;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final RichTextContent getContent() {
        return this.content;
    }

    @InterfaceC88464a
    public final RichTextEntity copy(@InterfaceC88465b Integer start_index, @InterfaceC88465b Integer end_index, @InterfaceC88465b RichTextContent content) {
        return new RichTextEntity(start_index, end_index, content);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof RichTextEntity)) {
            return false;
        }
        RichTextEntity richTextEntity = (RichTextEntity) other;
        return Intrinsics.m65267c(this.start_index, richTextEntity.start_index) && Intrinsics.m65267c(this.end_index, richTextEntity.end_index) && Intrinsics.m65267c(this.content, richTextEntity.content);
    }

    public int hashCode() {
        Integer num = this.start_index;
        int iHashCode = (num == null ? 0 : num.hashCode()) * 31;
        Integer num2 = this.end_index;
        int iHashCode2 = (iHashCode + (num2 == null ? 0 : num2.hashCode())) * 31;
        RichTextContent richTextContent = this.content;
        return iHashCode2 + (richTextContent != null ? richTextContent.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        Integer num = this.start_index;
        Integer num2 = this.end_index;
        RichTextContent richTextContent = this.content;
        StringBuilder sbM10485a = C8366s1.m10485a("RichTextEntity(start_index=", num, ", end_index=", num2, ", content=");
        sbM10485a.append(richTextContent);
        sbM10485a.append(Separators.RPAREN);
        return sbM10485a.toString();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}